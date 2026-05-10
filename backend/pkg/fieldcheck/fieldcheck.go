package fieldcheck

import "regexp"

var (
	cnMobileRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)
	trackingRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{5,63}$`)
	pickupCodeRe  = regexp.MustCompile(`^\d{6}$`)
)

func IsCNMobile(phone string) bool {
	return cnMobileRegex.MatchString(phone)
}

func IsTrackingNo(trackingNo string) bool {
	return trackingRegex.MatchString(trackingNo)
}

func IsPickupCode(code string) bool {
	return pickupCodeRe.MatchString(code)
}
