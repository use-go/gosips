package rtp

import (
	"errors"
	"fmt"
)

/** Represents an RTP RTPPacket.
 *  The RTPPacket class can be used to parse a RTPRawPacket instance if it represents RTP data.
 *  The class can also be used to create a new RTP packet according to the parameters specified by
 *  the user.
 */
type RTPPacket struct {
	receivetime *RTPTime
	header      *RTPHeader
	extension   *RTPExtension
	payload     []byte

	packet []byte
}

/** Creates an RTPPacket instance based upon the data in \c rawpack, optionally installing a memory manager.
 *  If successful, the data is moved from the raw packet to the RTPPacket instance.
 */
func NewRTPPacketFromRawPacket(rawpack *RawPacket) (*RTPPacket, error) {
	rp := &RTPPacket{}
	rp.receivetime = rawpack.GetReceiveTime().Clone()
	if err := rp.ParseRawPacket(rawpack); err != nil {
		return nil, err
	}
	return rp, nil
}

func (rp *RTPPacket) ParseRawPacket(rawpack *RawPacket) error {
	if !rawpack.IsRTP() { // If we didn't receive it on the RTP port, we'll ignore it
		return errors.New("is not RTP.")
	}

	rp.packet = make([]byte, len(rawpack.GetData()))
	copy(rp.packet, rawpack.GetData())

	rp.header = NewRTPHeader()
	if err := rp.header.Parse(rp.packet); err != nil {
		return err
	}

	// The version number should be correct
	if rp.header.version != RTP_VERSION {
		return fmt.Errorf("invalid version. %d vs %d", rp.header.version, RTP_VERSION)
	}

	// We'll check if rp is possibly a RTCP packet. For rp to be possible
	// the marker bit and payload type combined should be either an SR or RR
	// identifier
	if rp.header.marker != 0 {
		if rp.header.payloadtype == (RTP_RTCPTYPE_SR & 127) { // don't check high bit (rp was the marker!!)
			return errors.New("invalid payloadtype rtcp_sr.")
		}
		if rp.header.payloadtype == (RTP_RTCPTYPE_RR & 127) {
			return errors.New("invalid payloadtype rtcp_rr.")
		}
	}

	var numpadbytes, payloadoffset, payloadlength int

	payloadoffset = SIZEOF_RTPHEADER + 4*int(rp.header.csrccount)
	if rp.header.extension != 0 { // got header extension
		rp.extension = NewRTPExtension()
		if err := rp.extension.Parse(rp.packet[payloadoffset:]); err != nil {
			return err
		}
		payloadoffset += SIZEOF_RTPEXTENSION + 4*int(rp.extension.length)
	} else {
		rp.extension = nil
	}

	if rp.header.padding != 0 { // adjust payload length to take padding into account
		numpadbytes = int(rp.packet[len(rp.packet)-1]) // last byte contains number of padding bytes
		if numpadbytes > len(rp.packet)-payloadoffset {
			return errors.New("invalid padding.")
		}
	} else {
		numpadbytes = 0
	}

	payloadlength = len(rp.packet) - numpadbytes - payloadoffset
	if payloadlength < 0 {
		return errors.New("invalid payload length.")
	}

	return nil
}

/** Creates a new buffer for an RTP packet and fills in the fields according to the specified parameters.
 *  If \c maxpacksize is not equal to zero, an error is generated if the total packet size would exceed
 *  \c maxpacksize. The arguments of the constructor are self-explanatory. Note that the size of a header
 *  extension is specified in a number of 32-bit words. A memory manager can be installed.
 *  This constructor is similar to the other constructor, but here data is stored in an external buffer
 *  \c buffer with size \c buffersize. */
func NewPacket(payloadtype uint8,
	payloaddata []byte,
	seqnr uint16,
	timestamp uint32,
	ssrc uint32,
	gotmarker bool,
	numcsrcs uint8,
	csrcs []uint32,
	gotextension bool,
	extensionid uint16,
	extensionlen uint16,
	extensiondata []uint32) *RTPPacket {
	rp := &RTPPacket{}

	rp.receivetime = &RTPTime{0, 0}
	if err := rp.BuildPacket(payloadtype,
		payloaddata,
		seqnr,
		timestamp,
		ssrc,
		gotmarker,
		numcsrcs,
		csrcs,
		gotextension,
		extensionid,
		extensionlen,
		extensiondata); err != nil {
		return nil
	}

	return rp
}
func (rp *RTPPacket) BuildPacket(payloadtype uint8,
	payloaddata []byte,
	seqnr uint16,
	timestamp uint32,
	ssrc uint32,
	gotmarker bool,
	numcsrcs uint8,
	csrcs []uint32,
	gotextension bool,
	extensionid uint16,
	extensionlen uint16,
	extensiondata []uint32) error {
	if numcsrcs > RTP_MAXCSRCS {
		return errors.New("ERR_RTP_PACKET_TOOMANYCSRCS")
	}

	if payloadtype > 127 { // high bit should not be used
		return errors.New("ERR_RTP_PACKET_BADPAYLOADTYPE")
	}
	if payloadtype == 72 || payloadtype == 73 { // could cause confusion with rtcp types
		return errors.New("ERR_RTP_PACKET_BADPAYLOADTYPE")
	}

	var packetlength, packetoffset int
	packetlength = SIZEOF_RTPHEADER
	packetlength += int(numcsrcs) * 4 //sizeof(uint32_t)*((size_t)
	if gotextension {
		packetlength += SIZEOF_RTPEXTENSION   //(RTPExtensionHeader);
		packetlength += int(extensionlen) * 4 //sizeof(uint32_t)*((size_t)
	}
	packetlength += len(payloaddata) //payloadlen;
	rp.packet = make([]byte, packetlength)

	// Ok, now we'll just fill in...
	rp.header = NewRTPHeader()
	rp.header.version = RTP_VERSION
	rp.header.padding = 0
	if gotextension {
		rp.header.extension = 1
	} else {
		rp.header.extension = 0
	}
	rp.header.csrccount = numcsrcs
	if gotmarker {
		rp.header.marker = 1
	} else {
		rp.header.marker = 0
	}
	rp.header.payloadtype = payloadtype & 127
	rp.header.sequencenumber = seqnr
	rp.header.timestamp = timestamp
	rp.header.ssrc = ssrc
	if numcsrcs != 0 {
		rp.header.csrc = make([]uint32, numcsrcs)
		for i := uint8(0); i < numcsrcs; i++ {
			rp.header.csrc[i] = csrcs[i] //htonl(csrcs[i]);
		}
	}

	packetoffset = SIZEOF_RTPHEADER + int(numcsrcs)*4
	copy(rp.packet[0:packetoffset], rp.header.Encode())

	if gotextension {
		rp.extension = NewRTPExtension()
		rp.extension.id = extensionid
		rp.extension.length = extensionlen //sizeof(uint32_t);
		if extensionlen != 0 {
			rp.extension.data = make([]uint32, extensionlen)
			for i := uint16(0); i < extensionlen; i++ {
				rp.extension.data[i] = extensiondata[i]
			}
		}
		copy(rp.packet[packetoffset:packetoffset+SIZEOF_RTPEXTENSION+int(extensionlen)*4], rp.extension.Encode())

		packetoffset += SIZEOF_RTPEXTENSION + int(extensionlen)*4
	} else {
		rp.extension = nil
	}

	rp.payload = make([]byte, len(payloaddata))
	copy(rp.payload, payloaddata)
	copy(rp.packet[packetoffset:packetoffset+len(payloaddata)], payloaddata)

	return nil
}

/** Returns \c true if the RTP packet has a header extension and \c false otherwise. */
func (rp *RTPPacket) HasExtension() bool {
	return rp.header.extension != 0
}

/** Returns \c true if the marker bit was set and \c false otherwise. */
func (rp *RTPPacket) HasMarker() bool {
	return rp.header.marker != 0
}

/** Returns the number of CSRCs contained in rp packet. */
func (rp *RTPPacket) GetCSRCCount() uint8 {
	return rp.header.csrccount
}

/** Returns a specific CSRC identifier.
 *  Returns a specific CSRC identifier. The parameter \c num can go from 0 to GetCSRCCount()-1.
 */
func (rp *RTPPacket) GetCSRC(num uint8) uint32 {
	if num >= rp.header.csrccount {
		return 0
	}

	return rp.header.csrc[num]
}

/** Returns the payload type of the packet. */
func (rp *RTPPacket) GetPayloadType() uint8 {
	return rp.header.payloadtype
}

/** Returns the extended sequence number of the packet.
 *  Returns the extended sequence number of the packet. When the packet is just received,
 *  only the low $16$ bits will be set. The high 16 bits can be filled in later.
 */
// func (rp *RTPPacket) GetExtendedSequenceNumber() uint32 {
// 	return rp.extseqnr
// }

/** Returns the sequence number of rp packet. */
func (rp *RTPPacket) GetSequenceNumber() uint16 {
	return rp.header.sequencenumber //uint16(rp.extseqnr & 0x0000FFFF)
}

/** Sets the extended sequence number of rp packet to \c seq. */
// func (rp *RTPPacket) SetExtendedSequenceNumber(seq uint32) {
// 	rp.extseqnr = seq
// }

/** Returns the timestamp of rp packet. */
func (rp *RTPPacket) GetTimestamp() uint32 {
	return rp.header.timestamp
}

/** Returns the SSRC identifier stored in rp packet. */
func (rp *RTPPacket) GetSSRC() uint32 {
	return rp.header.ssrc
}

/** Returns a pointer to the actual payload data. */
func (rp *RTPPacket) GetPayload() []byte {
	return rp.payload
}

/** If a header extension is present, rp function returns the extension identifier. */
func (rp *RTPPacket) GetExtensionID() uint16 {
	return rp.extension.id
}

/** Returns the length of the header extension data. */
func (rp *RTPPacket) GetExtensionLength() uint16 {
	return rp.extension.length
}

/** Returns the header extension data. */
func (rp *RTPPacket) GetExtensionData() []uint32 {
	return rp.extension.data
}

/** Returns the time at which rp packet was received.
 *  When an RTPPacket instance is created from an RTPRawPacket instance, the raw packet's
 *  reception time is stored in the RTPPacket instance. This function then retrieves that
 *  time.
 */
func (rp *RTPPacket) GetReceiveTime() *RTPTime {
	return rp.receivetime
}

/** Returns a pointer to the data of the entire packet. */
func (rp *RTPPacket) GetPacket() []byte {
	return rp.packet
}

func (rp *RTPPacket) Dump() string {
	/*
		log.Printf("Payload type:                %d\n", rp.GetPayloadType())
		log.Printf("Sequence number:             0x%08x\n", rp.GetSequenceNumber())
		log.Printf("Timestamp:                   0x%08x\n", rp.GetTimestamp())
		log.Printf("SSRC:                        0x%08x\n", rp.GetSSRC())
		log.Printf("CSRC count:                  %d\n", rp.GetCSRCCount())
		for i := uint8(0); i < rp.GetCSRCCount(); i++ {
			log.Printf("    CSRC[%02d]:                0x%08x\n", i, rp.GetCSRC(i))
		}
		if rp.HasExtension() {
			log.Printf("    RTPExtension ID:            0x%04x\n", rp.GetExtensionID())
			log.Printf("    RTPExtension length:        %d\n", rp.GetExtensionLength())
		}
	*/
	return fmt.Sprintf("RTPPacket{\n\tMarker: %t\n\tPayloadType: %d\n\tSeq: %d\n\tTimestamp: %d\n\tSSRC: %d\n\tPayload length: %d\n}", rp.HasMarker(), rp.GetPayloadType(), rp.GetSequenceNumber(), rp.GetTimestamp(), rp.GetSSRC(), len(rp.GetPayload()))
}
