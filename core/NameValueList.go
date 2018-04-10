package core

import (
	"bytes"
	"container/list"
)

/**
* Implements a simple NameValue association with a quick lookup
* function (via a hash table) nameValueList class is not thread safe
* because it uses HashTables.
 */

type NameValueList struct {
	list.List

	indentation int
	listName    string // For debugging
	stringRep   string
	separator   string
}

func NewNameValueList(listName string) *NameValueList {
	listname := &NameValueList{}

	listname.listName = listName
	listname.separator = ";"

	return listname
}

func (nameValueList *NameValueList) GetIndentation() string {
	var retval bytes.Buffer
	for i := 0; i < nameValueList.indentation; i++ {
		retval.WriteString(" ")
	}
	return retval.String()
}

/*func (nameValueList *NameValueList) ConcatenateToTail(nvl *NameValueList) {
    nameValueList.Concatenate(objList, false)
}*/

func (nameValueList *NameValueList) Concatenate(nvl *NameValueList, topFlag bool) {
	if nvl == nil {
		return
	}

	if !topFlag {
		//nameValueList.PushBackList(nvl)
		for e := nvl.Front(); e != nil; e = e.Next() {
		    nameValueList.PushBack(e)
		}
	} else {
		//nameValueList.PushFrontList(nvl)
		// add given items to the end of the list.
		first := nameValueList.Front()
		for e := nvl.Front(); e != nil; e = e.Next() {
		    nameValueList.InsertBefore(e, first)
		}
	}
}

/**
 * string formatting function.
 */

func (nameValueList *NameValueList) Sprint(s string) {
	if s == "" {
		nameValueList.stringRep += nameValueList.GetIndentation()
		nameValueList.stringRep += "<null>\n"
		return
	}

	if s == "}" || s == "]" {
		nameValueList.indentation--
	}
	nameValueList.stringRep += nameValueList.GetIndentation()
	nameValueList.stringRep += s
	nameValueList.stringRep += "\n"
	if s == "{" || s == "[" {
		nameValueList.indentation++
	}
}

/**
         * Encode the list in semicolon separated form.
	 * @return an encoded string containing the objects in nameValueList list.
         * @since v1.0
*/
func (nameValueList *NameValueList) String() string {
	if nameValueList.Len() == 0 {
		return ""
	}

	var encoding bytes.Buffer //= new StringBuffer();
	for e := nameValueList.Front(); e != nil; e = e.Next() {
		nv := e.Value.(*NameValue)
		encoding.WriteString(nv.String())

		if e.Next() != nil {
			//println(nameValueList.separator);
			encoding.WriteString(nameValueList.separator)
		}

	}

	return encoding.String()
}

/**
 *  Set the separator (for encoding the list)
 * @since v1.0
 * @param sep is the new seperator (default is semicolon)
 */
func (nameValueList *NameValueList) SetSeparator(sep string) {
	nameValueList.separator = sep
}

func (nameValueList *NameValueList) AddNameValue(nv *NameValue) {
	if nv == nil {
		//throw new NullPointerException("null nv");
		return
	}
	nameValueList.PushBack(nv)
}

/**
* Add a name value record to nameValueList list.
 */
func (nameValueList *NameValueList) AddNameAndValue(name string, value interface{}) {
	nv := NewNameValue(name, value)
	nameValueList.AddNameValue(nv)
}

/**
* Set a namevalue object in nameValueList list.
 */
func (nameValueList *NameValueList) SetNameValue(nv *NameValue) {
	nameValueList.Delete(nv.name)
	nameValueList.AddNameValue(nv)
}

/**
* Set a namevalue object in nameValueList list.
 */
func (nameValueList *NameValueList) SetNameAndValue(name string, value interface{}) {
	nv := NewNameValue(name, value)
	nameValueList.SetNameValue(nv)
}

/**
         *  Compare if two NameValue lists are equal.
	 *@param otherObject  is the object to compare to.
	 *@return true if the two objects compare for equality.
*/
/*public boolean equals(Object otherObject) {
            if (!otherObject.getClass().equals
                (nameValueList.getClass())) {
                return false;
            }
            NameValueList other = (NameValueList) otherObject;

            if (nameValueList.size() != other.size()) {
		return false;
	    }
	    ListIterator li = nameValueList.listIterator();

	    while (li.hasNext()) {
		NameValue nv = (NameValue) li.next();
		boolean found = false;
	        ListIterator li1 = other.listIterator();
		while (li1.hasNext()) {
			NameValue nv1  = (NameValue) li1.next();
			// found a match so break;
			if (nv.equals(nv1))   {
			   found = true;
			   break;
			}
		}
		if (! found ) return false;
	    }
	    return true;
	}*/

/**
*  Do a lookup on a given name and return value associated with it.
 */
func (nameValueList *NameValueList) GetValue(name string) interface{} {
	nv := nameValueList.GetNameValue(name)
	if nv != nil {
		return nv.value
	}

	return nil
}

/**
* Get the NameValue record given a name.
* @since 1.0
 */
func (nameValueList *NameValueList) GetNameValue(name string) *NameValue {
	for e := nameValueList.Front(); e != nil; e = e.Next() {
		nv := e.Value.(*NameValue)
		if nv.GetName() == name {
			return nv
		}
	}

	return nil
}

/**
* Returns a boolean telling if nameValueList NameValueList
* has a record with nameValueList name
* @since 1.0
 */
func (nameValueList *NameValueList) HasNameValue(name string) bool {
	return nameValueList.GetNameValue(name) != nil
}

/**
* Remove the element corresponding to nameValueList name.
* @since 1.0
 */
func (nameValueList *NameValueList) Delete(name string) bool {
	for e := nameValueList.Front(); e != nil; e = e.Next() {
		nv := e.Value.(*NameValue)
		if nv.GetName() == name {
			nameValueList.Remove(e)
			return true
		}
	}

	return false
}

/**
 *Get a list of parameter names.
 *@return a list iterator that has the names of the parameters.
 */
func (nameValueList *NameValueList) GetNames() *list.List {
	ll := list.New()
	for e := nameValueList.Front(); e != nil; e = e.Next() {
		nv := e.Value.(*NameValue)
		ll.PushBack(nv.GetName())
	}
	return ll
}

func (nameValueList *NameValueList) Clone() interface{} {
	retval := &NameValueList{}
	retval.indentation = nameValueList.indentation
	retval.listName = nameValueList.listName
	retval.stringRep = nameValueList.stringRep
	retval.separator = nameValueList.separator

	li := list.New()
	for e := nameValueList.Front(); e != nil; e = e.Next() {
		nv := e.Value.(*NameValue)
		nnv := nv.Clone().(*NameValue)
		li.PushBack(nnv)
	}

	return retval
}

/** Get the parameter as a String.
 *@return the parameter as a string.
 */
func (nameValueList *NameValueList) GetParameter(name string) string {
	val := nameValueList.GetValue(name)
	if val == nil {
		return ""
	}

	return val.(string)
}
