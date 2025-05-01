package riot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/tyryoxan/API-RiftRadar/domain"
)

// Constants for HTTP client configuration
const (
	defaultTimeout  = 10 * time.Second
	maxIdleConns    = 10
	idleConnTimeout = 90 * time.Second

	// Rate limiting constants
	requestsPerSecond = 20
	requestsPerMinute = 100
	requestsPerHour   = 1000

	// Cache expiration
	cacheExpiration = 5 * time.Minute
)

type RiotPlayerRepo struct {
	apiKey string
	client *http.Client
	mu     sync.Mutex
	cache  map[string]cachedPlayer

	// Rate limiting
	requestTimes []time.Time
	rateLimitMu  sync.Mutex
}

type cachedPlayer struct {
	player domain.Player
	expiry time.Time
}

func NewRiotPlayerRepo() *RiotPlayerRepo {
	// Create a transport with connection pooling
	transport := &http.Transport{
		MaxIdleConns:        maxIdleConns,
		IdleConnTimeout:     idleConnTimeout,
		DisableCompression:  false,
		MaxIdleConnsPerHost: maxIdleConns,
	}

	// Create a client with the transport and timeout
	client := &http.Client{
		Transport: transport,
		Timeout:   defaultTimeout,
	}

	return &RiotPlayerRepo{
		apiKey:       os.Getenv("RIOT_API_KEY"),
		client:       client,
		cache:        make(map[string]cachedPlayer),
		requestTimes: make([]time.Time, 0, requestsPerHour),
	}
}

// checkRateLimit checks if we're within rate limits and waits if necessary
func (r *RiotPlayerRepo) checkRateLimit() {
	r.rateLimitMu.Lock()
	defer r.rateLimitMu.Unlock()

	now := time.Now()

	// Remove old requests from the tracking slice
	i := 0
	for i < len(r.requestTimes) && now.Sub(r.requestTimes[i]) > time.Hour {
		i++
	}
	if i > 0 {
		r.requestTimes = r.requestTimes[i:]
	}

	// Check hourly limit
	if len(r.requestTimes) >= requestsPerHour {
		// Wait until the oldest request is more than an hour old
		waitTime := time.Hour - now.Sub(r.requestTimes[0]) + time.Millisecond*100
		r.rateLimitMu.Unlock() // Unlock while waiting
		time.Sleep(waitTime)
		r.rateLimitMu.Lock()
		now = time.Now() // Update current time
	}

	// Check minute limit
	minuteOld := now.Add(-time.Minute)
	minuteCount := 0
	for i := len(r.requestTimes) - 1; i >= 0; i-- {
		if r.requestTimes[i].After(minuteOld) {
			minuteCount++
		} else {
			break
		}
	}

	if minuteCount >= requestsPerMinute {
		// Wait until we're under the minute limit
		waitTime := time.Minute - now.Sub(r.requestTimes[len(r.requestTimes)-requestsPerMinute]) + time.Millisecond*100
		r.rateLimitMu.Unlock() // Unlock while waiting
		time.Sleep(waitTime)
		r.rateLimitMu.Lock()
		now = time.Now() // Update current time
	}

	// Check second limit
	secondOld := now.Add(-time.Second)
	secondCount := 0
	for i := len(r.requestTimes) - 1; i >= 0; i-- {
		if r.requestTimes[i].After(secondOld) {
			secondCount++
		} else {
			break
		}
	}

	if secondCount >= requestsPerSecond {
		// Wait until we're under the second limit
		waitTime := time.Second - now.Sub(r.requestTimes[len(r.requestTimes)-requestsPerSecond]) + time.Millisecond*10
		r.rateLimitMu.Unlock() // Unlock while waiting
		time.Sleep(waitTime)
		r.rateLimitMu.Lock()
		now = time.Now() // Update current time
	}

	// Record this request
	r.requestTimes = append(r.requestTimes, now)
}

// Save is not implemented for the Riot API adapter
func (r *RiotPlayerRepo) Save(player domain.Player) error {
	return fmt.Errorf("save operation not supported for Riot API adapter")
}

// FindAll is not implemented for the Riot API adapter
func (r *RiotPlayerRepo) FindAll() []domain.Player {
	return []domain.Player{}
}

// GetPlayerByNameAndRegion fetches player information from the Riot API
func (r *RiotPlayerRepo) GetPlayerByNameAndRegion(name, region string) (domain.Player, error) {
	// Convert region to lowercase for consistency
	region = strings.ToLower(region)

	// Check cache first
	cacheKey := fmt.Sprintf("%s:%s", name, region)
	if player, found := r.getCachedPlayer(cacheKey); found {
		return player, nil
	}

	// Map region to Riot API region code
	regionCode := r.mapRegionToCode(region)

	// For backward compatibility, use the region as the tag
	player, err := r.getPlayerInfo(name, regionCode, region)
	if err != nil {
		return domain.Player{}, err
	}

	// Cache the result
	r.cachePlayer(cacheKey, player)
	return player, nil
}

// GetPlayerByNameAndTag fetches player information from the Riot API using name and tag
func (r *RiotPlayerRepo) GetPlayerByNameAndTag(name, tag string) (domain.Player, error) {
	// Default to EUW region for now
	region := "euw"

	// Check cache first
	cacheKey := fmt.Sprintf("%s:%s", name, tag)
	if player, found := r.getCachedPlayer(cacheKey); found {
		return player, nil
	}

	player, err := r.getPlayerInfo(name, tag, region)
	if err != nil {
		return domain.Player{}, err
	}

	// Cache the result
	r.cachePlayer(cacheKey, player)
	return player, nil
}

// getCachedPlayer retrieves a player from the cache if it exists and is not expired
func (r *RiotPlayerRepo) getCachedPlayer(key string) (domain.Player, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cached, exists := r.cache[key]
	if !exists {
		return domain.Player{}, false
	}

	// Check if the cached data is expired
	if time.Now().After(cached.expiry) {
		delete(r.cache, key)
		return domain.Player{}, false
	}

	return cached.player, true
}

// cachePlayer adds a player to the cache with an expiry time
func (r *RiotPlayerRepo) cachePlayer(key string, player domain.Player) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Cache for 5 minutes
	r.cache[key] = cachedPlayer{
		player: player,
		expiry: time.Now().Add(5 * time.Minute),
	}
}

// getPlayerInfo is a helper method to fetch player information
func (r *RiotPlayerRepo) getPlayerInfo(name, tag, region string) (domain.Player, error) {
	// Step 1: Get summoner information
	summoner, err := r.getSummonerByName(name, tag)
	if err != nil {
		return domain.Player{}, err
	}

	// Map region to Riot API region code for league info
	regionCode := r.mapRegionToCode(region)

	// Step 2: Get league information
	leagueInfo, err := r.getLeagueInfoBySummonerId(summoner.PUUID, regionCode)
	if err != nil {
		return domain.Player{}, err
	}

	// Create player object from the API response
	player := domain.Player{
		Name:   summoner.Name,
		Region: region,
		Elo:    leagueInfo.Tier,
		Lp:     leagueInfo.LeaguePoints,
	}

	return player, nil
}

// RiotAccount represents the response from the Riot Account API
type RiotAccount struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

// Summoner represents the response from the Riot API summoner endpoint
type Summoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

// LeagueEntry represents the response from the Riot API league endpoint
type LeagueEntry struct {
	LeagueID     string `json:"leagueId"`
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	HotStreak    bool   `json:"hotStreak"`
	Veteran      bool   `json:"veteran"`
	FreshBlood   bool   `json:"freshBlood"`
	Inactive     bool   `json:"inactive"`
}

// Constants for API endpoints
const (
	riotAccountAPI = "https://europe.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s"
	riotLeagueAPI  = "https://%s.api.riotgames.com/lol/league/v4/entries/by-puuid/%s"
)

// getSummonerByName fetches summoner information from the Riot API
func (r *RiotPlayerRepo) getSummonerByName(name, tag string) (Summoner, error) {
	// Apply rate limiting
	r.checkRateLimit()

	url := fmt.Sprintf(riotAccountAPI, name, tag)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Summoner{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-Riot-Token", r.apiKey)

	// Use the shared HTTP client
	resp, err := r.client.Do(req)
	if err != nil {
		return Summoner{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		// Handle rate limiting from Riot API
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			seconds := 1 // Default to 1 second if parsing fails
			fmt.Sscanf(retryAfter, "%d", &seconds)
			time.Sleep(time.Duration(seconds) * time.Second)
			// Retry the request
			return r.getSummonerByName(name, tag)
		}
	} else if resp.StatusCode != http.StatusOK {
		return Summoner{}, fmt.Errorf("riot API returned status code %d for %s", resp.StatusCode, url)
	}

	var account RiotAccount
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return Summoner{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Create a Summoner object from the RiotAccount
	summoner := Summoner{
		PUUID: account.PUUID,
		Name:  account.GameName, // Use GameName as the Name field
	}

	return summoner, nil
}

// getLeagueInfoBySummonerId fetches league information from the Riot API
func (r *RiotPlayerRepo) getLeagueInfoBySummonerId(summonerId, region string) (LeagueEntry, error) {
	// Apply rate limiting
	r.checkRateLimit()

	url := fmt.Sprintf(riotLeagueAPI, region, summonerId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LeagueEntry{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-Riot-Token", r.apiKey)

	// Use the shared HTTP client
	resp, err := r.client.Do(req)
	if err != nil {
		return LeagueEntry{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		// Handle rate limiting from Riot API
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			seconds := 1 // Default to 1 second if parsing fails
			fmt.Sscanf(retryAfter, "%d", &seconds)
			time.Sleep(time.Duration(seconds) * time.Second)
			// Retry the request
			return r.getLeagueInfoBySummonerId(summonerId, region)
		}
	} else if resp.StatusCode != http.StatusOK {
		return LeagueEntry{}, fmt.Errorf("riot API returned status code %d for %s", resp.StatusCode, url)
	}

	var leagueEntries []LeagueEntry
	if err := json.NewDecoder(resp.Body).Decode(&leagueEntries); err != nil {
		return LeagueEntry{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Find the solo queue entry
	for _, entry := range leagueEntries {
		if entry.QueueType == "RANKED_SOLO_5x5" {
			return entry, nil
		}
	}

	// If no solo queue entry is found, return a default entry with "UNRANKED" tier and 0 LP
	return LeagueEntry{
		Tier:         "UNRANKED",
		LeaguePoints: 0,
	}, nil
}

// mapRegionToCode maps a region name to a Riot API region code
func (r *RiotPlayerRepo) mapRegionToCode(region string) string {
	regionMap := map[string]string{
		"euw":  "euw1",
		"eune": "eun1",
		"na":   "na1",
		"kr":   "kr",
		"jp":   "jp1",
		"br":   "br1",
		"lan":  "la1",
		"las":  "la2",
		"oce":  "oc1",
		"tr":   "tr1",
		"ru":   "ru",
	}

	code, ok := regionMap[region]
	if !ok {
		return "euw1" // Default to EUW if region not found
	}

	return code
}
