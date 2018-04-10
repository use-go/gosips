package header

import (
	"github.com/use-go/gosips/core"
)

/**
* ContentLanguage list of headers. (Should this be a list?)
 */
type ContentLanguageList struct {
	SIPHeaderList
}

/** Default constructor
 */
func NewContentLanguageList() *ContentLanguageList {
	this := &ContentLanguageList{}
	this.SIPHeaderList.super(core.SIPHeaderNames_CONTENT_LANGUAGE)
	return this
}
