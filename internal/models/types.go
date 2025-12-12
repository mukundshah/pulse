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

type IPVersionType string

const (
	IPVersionTypeIPv4 IPVersionType = "ipv4"
	IPVersionTypeIPv6 IPVersionType = "ipv6"
)

type RetryType string

const (
	RetryTypeNone        RetryType = "none"
	RetryTypeFixed       RetryType = "fixed"
	RetryTypeLinear      RetryType = "linear"
	RetryTypeExponential RetryType = "exponential"
)

type RetryJitterType string

const (
	RetryJitterTypeNone         RetryJitterType = "none"
	RetryJitterTypeFull         RetryJitterType = "full"
	RetryJitterTypeEqual        RetryJitterType = "equal"
	RetryJitterTypeDecorrelated RetryJitterType = "decorrelated"
)

type UnitType string

const (
	UnitTypeMs UnitType = "ms"
	UnitTypeS  UnitType = "s"
)

type DNSRecordType string

const (
	DNSRecordTypeA     DNSRecordType = "A"
	DNSRecordTypeAAAA  DNSRecordType = "AAAA"
	DNSRecordTypeCNAME DNSRecordType = "CNAME"
	DNSRecordTypeMX    DNSRecordType = "MX"
	DNSRecordTypeNS    DNSRecordType = "NS"
	DNSRecordTypeSOA   DNSRecordType = "SOA"
	DNSRecordTypeSRV   DNSRecordType = "SRV"
	DNSRecordTypeTXT   DNSRecordType = "TXT"
)

type DNSResolverProtocolType string

const (
	DNSResolverProtocolUDP DNSResolverProtocolType = "udp"
	DNSResolverProtocolTCP DNSResolverProtocolType = "tcp"
)
