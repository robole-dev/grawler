package grawl

import (
	"github.com/gocolly/colly/v2"
	"net/http"
	"strings"
)

func StatusAbbreviation(code int) string {
	switch code {
	case http.StatusContinue:
		return "CT"
	case http.StatusSwitchingProtocols:
		return "SP"
	case http.StatusProcessing:
		return "PR"
	case http.StatusEarlyHints:
		return "EH"
	case http.StatusOK:
		return "OK"
	case http.StatusCreated:
		return "CR"
	case http.StatusAccepted:
		return "AC"
	case http.StatusNonAuthoritativeInfo:
		return "NA"
	case http.StatusNoContent:
		return "NC"
	case http.StatusResetContent:
		return "RC"
	case http.StatusPartialContent:
		return "PC"
	case http.StatusMultiStatus:
		return "MS"
	case http.StatusAlreadyReported:
		return "AR"
	case http.StatusIMUsed:
		return "IU"
	case http.StatusMultipleChoices:
		return "MC"
	case http.StatusMovedPermanently:
		return "MP"
	case http.StatusFound:
		return "FD"
	case http.StatusSeeOther:
		return "SO"
	case http.StatusNotModified:
		return "NM"
	case http.StatusUseProxy:
		return "UP"
	case http.StatusTemporaryRedirect:
		return "TR"
	case http.StatusPermanentRedirect:
		return "PR"
	case http.StatusBadRequest:
		return "BR"
	case http.StatusUnauthorized:
		return "UN"
	case http.StatusPaymentRequired:
		return "PR"
	case http.StatusForbidden:
		return "FB"
	case http.StatusNotFound:
		return "NF"
	case http.StatusMethodNotAllowed:
		return "MN"
	case http.StatusNotAcceptable:
		return "NA"
	case http.StatusProxyAuthRequired:
		return "PA"
	case http.StatusRequestTimeout:
		return "RT"
	case http.StatusConflict:
		return "CF"
	case http.StatusGone:
		return "GN"
	case http.StatusLengthRequired:
		return "LR"
	case http.StatusPreconditionFailed:
		return "PF"
	case http.StatusRequestEntityTooLarge:
		return "RL"
	case http.StatusRequestURITooLong:
		return "RL"
	case http.StatusUnsupportedMediaType:
		return "UM"
	case http.StatusRequestedRangeNotSatisfiable:
		return "RR"
	case http.StatusExpectationFailed:
		return "EF"
	case http.StatusTeapot:
		return "TP"
	case http.StatusMisdirectedRequest:
		return "MR"
	case http.StatusUnprocessableEntity:
		return "UE"
	case http.StatusLocked:
		return "LK"
	case http.StatusFailedDependency:
		return "FD"
	case http.StatusTooEarly:
		return "TE"
	case http.StatusUpgradeRequired:
		return "UR"
	case http.StatusPreconditionRequired:
		return "PR"
	case http.StatusTooManyRequests:
		return "TR"
	case http.StatusRequestHeaderFieldsTooLarge:
		return "RF"
	case http.StatusUnavailableForLegalReasons:
		return "UL"
	case http.StatusInternalServerError:
		return "IE"
	case http.StatusNotImplemented:
		return "NI"
	case http.StatusBadGateway:
		return "BG"
	case http.StatusServiceUnavailable:
		return "SU"
	case http.StatusGatewayTimeout:
		return "GT"
	case http.StatusHTTPVersionNotSupported:
		return "HS"
	case http.StatusVariantAlsoNegotiates:
		return "VN"
	case http.StatusInsufficientStorage:
		return "IS"
	case http.StatusLoopDetected:
		return "LD"
	case http.StatusNotExtended:
		return "NE"
	case http.StatusNetworkAuthenticationRequired:
		return "NR"
	default:
		return "??"
	}
}

func IsXmlResponse(resp *colly.Response) bool {
	contentType := strings.ToLower(resp.Headers.Get("Content-Type"))
	isXMLFile := strings.HasSuffix(strings.ToLower(resp.Request.URL.Path), ".xml") || strings.HasSuffix(strings.ToLower(resp.Request.URL.Path), ".xml.gz")
	isXmlContentType := strings.Contains(contentType, "xml")
	isHtmlContentType := strings.Contains(contentType, "html")

	return !isHtmlContentType && (isXMLFile || isXmlContentType)
}

func IsHtmlResponse(resp *colly.Response) bool {
	contentType := strings.ToLower(resp.Headers.Get("Content-Type"))
	return strings.Contains(contentType, "html")
}
