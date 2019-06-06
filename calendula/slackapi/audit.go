package slackapi

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	misc "github.com/xinnige/asteraceae/calendula/astermisc"
)

const (
	// APIURL defines for slack web api endpoint
	APIURL = "https://slack.com/api/"
	// AUDITURL defines for slack audit web api endpoint
	AUDITURL = "https://api.slack.com/audit/v1/"

	ctypeJSON             = "application/json"
	maxLimit              = 9999
	errPaginationComplete = errorString("pagination complete")
)

// AuditLogsOption provided when getting audit logs.
type AuditLogsOption func(*AuditLogPagination)

// AuditLogsOptionLatest sets
// Unix timestamp of the most recent audit event to include (inclusive).
// see https://api.slack.com/docs/audit-logs-api
func AuditLogsOptionLatest(n int) AuditLogsOption {
	return func(p *AuditLogPagination) {
		p.latest = n
	}
}

// AuditLogsOptionOldest sets
// Unix timestamp of the least recent audit event to include (inclusive).
// see https://api.slack.com/docs/audit-logs-api
func AuditLogsOptionOldest(n int) AuditLogsOption {
	return func(p *AuditLogPagination) {
		p.oldest = n
	}
}

// AuditLogsOptionLimit sets
// Number of results to optimistically return, maximum 9999
// see https://api.slack.com/docs/audit-logs-api
func AuditLogsOptionLimit(n int) AuditLogsOption {
	return func(p *AuditLogPagination) {
		p.limit = n
	}
}

// AuditLogsOptionAction filters by	Name of the action
// see https://api.slack.com/docs/audit-logs-api
func AuditLogsOptionAction(action string) AuditLogsOption {
	return func(p *AuditLogPagination) {
		p.action = action
	}
}

// AuditLogsOptionActor filters by User ID who initiated the action.
// see https://api.slack.com/docs/audit-logs-api
func AuditLogsOptionActor(actor string) AuditLogsOption {
	return func(p *AuditLogPagination) {
		p.actor = actor
	}
}

// AuditLogsOptionEntity filters by ID of the target entity of the action
// see https://api.slack.com/docs/audit-logs-api
func AuditLogsOptionEntity(entity string) AuditLogsOption {
	return func(p *AuditLogPagination) {
		p.entity = entity
	}
}

type auditlogResponseFull struct {
	Entries  []AuditEntry     `json:"entries,omitempty"`
	Metadata ResponseMetadata `json:"response_metadata"`
}

// AuditEntry contains info of an entity
type AuditEntry struct {
	ID         string       `json:"id"`
	DateCreate json.Number  `json:"date_create"`
	Action     string       `json:"action"`
	Actor      AuditActor   `json:"actor"`
	Entity     AuditEntity  `json:"entity"`
	Context    AuditContext `json:"context"`
}

// AuditActor contains info of an actor
type AuditActor struct {
	Type string    `json:"type"`
	User AuditUser `json:"user"`
}

// AuditEntity contains info of an entity
type AuditEntity struct {
	Type string    `json:"type"`
	User AuditUser `json:"user"`
}

// AuditUser contains info of a user
type AuditUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// AuditContext contains info of a context
type AuditContext struct {
	UserAgent string        `json:"ua"`
	IPAddress string        `json:"ip_address"`
	Location  AuditLocation `json:"location"`
}

// AuditLocation contains info of a location
type AuditLocation struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

// AuditLogPagination allows for paginating over the audit logs
type AuditLogPagination struct {
	Entries      []AuditEntry
	limit        int
	latest       int
	oldest       int
	action       string
	entity       string
	actor        string
	previousResp *ResponseMetadata
	c            *Client
	values       url.Values
}

func newAuditLogPagination(c *Client, options ...AuditLogsOption) (p AuditLogPagination) {
	p = AuditLogPagination{
		c:     c,
		limit: maxLimit, // per slack api documentation.
	}

	for _, opt := range options {
		opt(&p)
	}
	return p
}

// Done checks if the pagination has completed
func (AuditLogPagination) Done(err error) bool {
	return err == errPaginationComplete
}

// Failure checks if pagination failed.
func (p AuditLogPagination) Failure(err error) error {
	if p.Done(err) {
		return nil
	}

	return err
}

func auditlogRequest(ctx context.Context, client *Client, path, token string, values url.Values) (*auditlogResponseFull, error) {
	response := &auditlogResponseFull{}
	err := misc.GetJSON(ctx, client.client, AUDITURL+path, token, values, response, client.method, client)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p AuditLogPagination) setValues(values *url.Values) {
	if p.latest != 0 {
		values.Add("latest", strconv.Itoa(p.latest))
	}
	if p.oldest != 0 {
		values.Add("oldest", strconv.Itoa(p.oldest))
	}
	if p.action != "" {
		values.Add("action", p.action)
	}
	if p.actor != "" {
		values.Add("actor", p.actor)
	}
	if p.entity != "" {
		values.Add("entity", p.entity)
	}

}

// Next iters paging of audit logs
func (p AuditLogPagination) Next(ctx context.Context) (_ AuditLogPagination, err error) {
	var (
		resp *auditlogResponseFull
	)

	if p.c == nil || (p.previousResp != nil && p.previousResp.Cursor == "") {
		return p, errPaginationComplete
	}

	p.previousResp = p.previousResp.initialize()

	values := url.Values{
		"limit":  {strconv.Itoa(p.limit)},
		"cursor": {p.previousResp.Cursor},
	}
	p.setValues(&values)

	if resp, err = auditlogRequest(ctx, p.c, "logs", p.c.token, values); err != nil {
		return p, err
	}

	p.c.Debugf("GetAuditLogs: got %d entries; metadata %v", len(resp.Entries), resp.Metadata)
	p.Entries = resp.Entries
	p.previousResp = &resp.Metadata

	return p, nil
}

// GetAuditLogsPaginated unarchives the given channel
// see https://api.slack.com/methods/channels.unarchive
func (client *Client) GetAuditLogsPaginated(options ...AuditLogsOption) AuditLogPagination {
	return newAuditLogPagination(client, options...)
}

// ListAuditLogs fetches logs in a paginated fashion, see GetAuditLogsPaginated for usage.
func (client *Client) ListAuditLogs(limit, latest, oldest int, action, actor, entity string) (entries []AuditEntry, err error) {
	opts := make([]AuditLogsOption, 0)
	if limit != 0 {
		opts = append(opts, AuditLogsOptionLimit(limit))
	}
	if latest != 0 {
		opts = append(opts, AuditLogsOptionLatest(latest))
	}
	if oldest != 0 {
		opts = append(opts, AuditLogsOptionOldest(oldest))
	}
	if actor != "" {
		opts = append(opts, AuditLogsOptionActor(actor))
	}
	if action != "" {
		opts = append(opts, AuditLogsOptionAction(action))
	}
	if entity != "" {
		opts = append(opts, AuditLogsOptionEntity(entity))
	}
	p := newAuditLogPagination(client, opts...)
	ctx := context.Background()
	results := make([]AuditEntry, 0)

	for ; !p.Done(err); p, err = p.Next(ctx) {
		results = append(results, p.Entries...)
	}
	return results, p.Failure(err)
}
