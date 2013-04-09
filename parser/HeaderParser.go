package parser

import (
	"strings"
	"gosip/core"
	"gosip/header"
)

type HeaderParser interface{
	Parse()  (sh header.SIPHeader, ParseException error); 
}

/** Generic header parser class. The parsers for various headers extend this
* class. To create a parser for a new header, extend this class and change
* the createParser class.
*/

type HeaderParserImpl struct{
 	ParserImpl
}    
	 
    /** Creates new HeaderParser
     * @param String to parse.
     */
    func NewHeaderParserImpl(header string) *HeaderParserImpl {
    	this := &HeaderParserImpl{}
    	
    	this.ParserImpl.super(header);
    	this.ParserImpl.GetLexer().SetLexerName("command_keywordLexer");
        
        return this;
    }
    
    func NewHeaderParserImplFromLexer(lexer core.Lexer) *HeaderParserImpl {
    	this := &HeaderParserImpl{}
    	
        this.SetLexer(lexer);
    	this.ParserImpl.GetLexer().SetLexerName("command_keywordLexer");
    	
    	return this;
    }
    
    func (this *HeaderParserImpl) super(header string) {
    	this.ParserImpl.super(header);
    	this.ParserImpl.GetLexer().SetLexerName("command_keywordLexer");
    }
    
    func (this *HeaderParserImpl) superFromLexer(lexer core.Lexer) {
    	this.SetLexer(lexer);
    	this.ParserImpl.GetLexer().SetLexerName("command_keywordLexer");
    }
    
    /** Parse the weekday field
     * @return an integer with the calendar content for wkday.
     */
    func (this *HeaderParserImpl) Wkday()  (wk int, ParseException error){
		this.Dbg_enter("wkday");
		defer this.Dbg_leave("wkday");
	//try {
        	tok := this.GetLexer().Ttoken();
        	id := strings.ToLower(tok);
	
        	if strings.ToLower(core.SIPDateNames_MON) == (id) { 
				return core.SIPCalendar_MONDAY, nil;
			}else if strings.ToLower(core.SIPDateNames_TUE)==(id){
				return core.SIPCalendar_TUESDAY, nil;
        	}else if strings.ToLower(core.SIPDateNames_WED)==(id){
				return core.SIPCalendar_WEDNESDAY, nil;
        	}else if strings.ToLower(core.SIPDateNames_THU)==(id){ 
				return core.SIPCalendar_THURSDAY, nil;
        	}else if strings.ToLower(core.SIPDateNames_FRI)==(id){ 
				return core.SIPCalendar_FRIDAY, nil;
        	}else if strings.ToLower(core.SIPDateNames_SAT)==(id){ 
				return core.SIPCalendar_SATURDAY, nil;
        	}else if strings.ToLower(core.SIPDateNames_SUN)==(id){ 
				return core.SIPCalendar_SUNDAY, nil;
        	}else{
        	  return -1, this.CreateParseException("bad wkday" );
        	}
	//} finally {
	//	dbg_leave("wkday");
	//}
        
    }
    
    /** parse and return a date field.
     *@return a date structure with the parsed value.
     */
    /*protected Calendar date() throws ParseException {
        try  {
            Calendar retval =
            Calendar.getInstance(TimeZone.getTimeZone("GMT"));
            String s1 = lexer.number();
            int day = Integer.parseInt(s1);
            if (day <= 0 || day >= 31)
                throw createParseException("Bad day ");
            retval.set(Calendar.DAY_OF_MONTH,day);
            lexer.match(' ');
            String month = lexer.ttoken().toLowerCase();
            if (month.equals("jan"))  {
                retval.set(Calendar.MONTH,Calendar.JANUARY);
            } else if (month.equals("feb")) {
                retval.set(Calendar.MONTH,Calendar.FEBRUARY);
            } else if (month.equals("mar")) {
                retval.set(Calendar.MONTH,Calendar.MARCH);
            } else if (month.equals("apr")) {
                retval.set(Calendar.MONTH,Calendar.APRIL);
            } else if (month.equals("may")) {
                retval.set(Calendar.MONTH,Calendar.MAY);
            } else if (month.equals("jun")) {
                retval.set(Calendar.MONTH,Calendar.JUNE);
            } else if (month.equals("jul")) {
                retval.set(Calendar.MONTH,Calendar.JULY);
            } else if (month.equals("aug")) {
                retval.set(Calendar.MONTH,Calendar.AUGUST);
            } else if (month.equals("sep")) {
                retval.set(Calendar.MONTH,Calendar.SEPTEMBER);
            } else if (month.equals("oct")) {
                retval.set(Calendar.MONTH,Calendar.OCTOBER);
            } else if (month.equals("nov")) {
                retval.set(Calendar.MONTH,Calendar.NOVEMBER);
            } else if (month.equals("dec")) {
                retval.set(Calendar.MONTH,Calendar.DECEMBER);
            }
            lexer.match(' ');
            String s2 = lexer.number();
            int yr = Integer.parseInt(s2);
            retval.set(Calendar.YEAR,yr);
            return retval;
            
        } catch (Exception ex) {
            throw createParseException("bad date field" );
        }
        
    }*/
    
    /** Set the time field. This has the format hour:minute:second
     */
    /*protected void time(Calendar calendar) throws ParseException {
        try {
            String s = lexer.number();
            int hour = Integer.parseInt(s);
            calendar.set(Calendar.HOUR_OF_DAY,hour);
            lexer.match(':');
            s = lexer.number();
            int min = Integer.parseInt(s);
            calendar.set(Calendar.MINUTE,min);
            lexer.match(':');
            s = lexer.number();
            int sec = Integer.parseInt(s);
            calendar.set(Calendar.SECOND,sec);
        } catch (Exception ex) {
            throw createParseException ("error processing time " );
            
        }
        
    }*/
   
    
    /** Parse the SIP header from the buffer and return a parsed
     * structure.
     *@throws ParseException if there was an error parsing.
     */
    func (this *HeaderParserImpl) Parse() (sh header.SIPHeader, ParseException error) {
            lexer:= this.GetLexer();
            
            name := lexer.GetNextTokenByDelim(':');
            lexer.ConsumeK(1);
            body := strings.TrimSpace(lexer.GetLine());
            // we dont set any fields because the header is
            // ok
            retval := header.NewExtensionHeaderImpl(name);
            retval.SetValue(body);
            return retval, nil;
            
    }
 
    /** Parse the header name until the colon  and chew WS after that.
     */
    func (this *HeaderParserImpl) HeaderName(tok int) {
            this.GetLexer().Match(tok);
            this.GetLexer().SPorHT();
            this.GetLexer().Match(':');
            this.GetLexer().SPorHT();
    }
	