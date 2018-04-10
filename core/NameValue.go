package core

type NameValue struct {
	isQuotedString bool
	separator      string
	quotes         string
	name           string
	value          interface{} //only accept nil and string
}

func NewNameValue(n string, v interface{}) *NameValue {
	if v != nil {
		if _, ok := v.(string); !ok {
			panic("value must be nil or string type")
		}
	}

	nameValue := &NameValue{}

	nameValue.name = n
	nameValue.value = v
	nameValue.separator = SIPSeparatorNames_EQUALS
	nameValue.quotes = ""

	return nameValue
}

/**
* Set the separator for the encoding method below.
 */
func (nameValue *NameValue) SetSeparator(sep string) {
	nameValue.separator = sep
}

/** A flag that indicates that doublequotes should be put around the
* value when encoded
*(for example name=value when value is doublequoted).
 */
func (nameValue *NameValue) SetQuotedValue() {
	nameValue.isQuotedString = true
	nameValue.quotes = SIPSeparatorNames_DOUBLE_QUOTE
}

/** Return true if the value is quoted in doublequotes.
 */
func (nameValue *NameValue) IsValueQuoted() bool {
	return nameValue.isQuotedString
}

func (nameValue *NameValue) GetName() string {
	return nameValue.name
}

func (nameValue *NameValue) GetValue() interface{} {
	return nameValue.value
}

/**
* Set the name member
 */
func (nameValue *NameValue) SetName(n string) {
	nameValue.name = n
}

/**
* Set the value member
 */
func (nameValue *NameValue) SetValue(v interface{}) {
	if v != nil {
		if _, ok := v.(string); !ok {
			panic("value must be nil or string type")
		}
	}
	nameValue.value = v
}

/**
	* Get the encoded representation of nameValue namevalue object.
        * Added doublequote for encoding doublequoted values
	* (bug reported by Kirby Kiem).
	*@since 1.0
	*@return an encoded name value (eg. name=value) string.
*/
func (nameValue *NameValue) String() string {
	if nameValue.name != "" && nameValue.value != nil {
		return nameValue.name + nameValue.separator + nameValue.quotes + nameValue.value.(string) + nameValue.quotes
	} else if nameValue.name == "" && nameValue.value != nil {
		return nameValue.quotes + nameValue.value.(string) + nameValue.quotes
	} else if nameValue.name != "" && nameValue.value == nil {
		return nameValue.name
	} else {
		return ""
	}
}

func (nameValue *NameValue) Clone() interface{} {
	retval := &NameValue{}
	retval.separator = nameValue.separator
	retval.isQuotedString = nameValue.isQuotedString
	retval.quotes = nameValue.quotes
	retval.name = nameValue.name
	if nameValue.value != nil {
		retval.value = nameValue.value.(string)
	}
	return retval
}

/**
* Equality comparison predicate.
 */
/*public boolean equals( Object other) {
	if (! other.getClass().equals(nameValue.getClass()))  return false;
        NameValue that = (NameValue) other;
	if (nameValue == that) return true;
	if (nameValue.name  == null && that.name != null ||
	   nameValue.name != null && that.name == null) return false;
	if (nameValue.name != null && that.name != null &&
		nameValue.name.compareToIgnoreCase(that.name) != 0)
			return false;
	if ( nameValue.value != null && that.value == null ||
	     nameValue.value == null && that.value != null) return false;
	if (nameValue.value == that.value) return true;
	if (value instanceof String) {
		// Quoted string comparisions are case sensitive.
	        if (isQuotedString)
			return nameValue.value.equals(that.value);
		String val = (String) nameValue.value;
		String val1 = (String) that.value;
		return val.compareToIgnoreCase(val1) == 0 ;
	} else return nameValue.value.equals(that.value);
}*/
