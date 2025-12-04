package models

type CheckType string

const (
	CheckTypeHTTP      CheckType = "http"
	CheckTypeTCP       CheckType = "tcp"
	CheckTypeDNS       CheckType = "dns"
	CheckTypeBrowser   CheckType = "browser"
	CheckTypeHeartbeat CheckType = "heartbeat"
)

type CheckRunStatus string

const (
	CheckRunStatusSuccess CheckRunStatus = "success"
	CheckRunStatusFail    CheckRunStatus = "fail"
	CheckRunStatusTimeout CheckRunStatus = "timeout"
	CheckRunStatusError   CheckRunStatus = "error"
	CheckRunStatusUnknown CheckRunStatus = "unknown"
)
