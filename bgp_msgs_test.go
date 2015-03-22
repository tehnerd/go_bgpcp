package bgp

import (
	"encoding/hex"
	"fmt"
	"testing"
)

const (
	hexOpenMsg               = "ffffffffffffffffffffffffffffffff003b0104fde8005a0a0000021e02060104000100010202800002020200020440020078020641040000fde8"
	hexUpdate1               = "ffffffffffffffffffffffffffffffff00360200000015400101004002004003040a0000024005040000006420010101012001010102"
	hexUpdate2               = "ffffffffffffffffffffffffffffffff0038020000001c400101004002004003040a00000280040400000078400504000000642001010103"
	hexUpdate3               = "ffffffffffffffffffffffffffffffff00170200000000"
	hexUpdate4               = "ffffffffffffffffffffffffffffffff005102000000364001010240021a02060000000100000002000000030000000400000005000000064003040a00000240050400000064c00804ffff000118010b01"
	hexKA                    = "ffffffffffffffffffffffffffffffff001304"
	hexNotification          = "ffffffffffffffffffffffffffffffff0015030607"
	hexIPv6NLRI              = "302a00bdc0e003"
	hexIPv6_MP_REACH         = "00020110200107f800200101000000000245018000302a00bdc0e003"
	hexIPv6_MP_REACH_NLRI_PA = "900e001c00020110200107f800200101000000000245018000302a00bdc0e003"
)

func TestDecodeMsgHeader(t *testing.T) {
	encodedOpen, _ := hex.DecodeString(hexOpenMsg)
	_, err := DecodeMsgHeader(encodedOpen)
	if err != nil {
		fmt.Println(err)
		t.Errorf("error during bgp msg header decoding")
	}
}

func TestEncodeMsgHeader(t *testing.T) {
	encodedOpen, _ := hex.DecodeString(hexOpenMsg)
	msgHdr := MsgHeader{Length: 59, Type: 1}
	encMsgHdr, err := EncodeMsgHeader(&msgHdr)
	if err != nil {
		fmt.Println(err)
		t.Errorf("error during bgp msg header encoding")
	}
	if len(encMsgHdr) != 19 {
		fmt.Println(len(encMsgHdr))
		fmt.Println(encMsgHdr)
		t.Errorf("error in len of encoded hdr")
	}
	for cntr := 0; cntr < len(encMsgHdr); cntr++ {
		if encMsgHdr[cntr] != encodedOpen[cntr] {
			t.Errorf("byte of encoded msg is not equal to etalon's msg")
		}
	}
}

func TestDecodeOpenMsg(t *testing.T) {
	encodedOpen, _ := hex.DecodeString(hexOpenMsg)
	openMsg, err := DecodeOpenMsg(encodedOpen[19:])
	if err != nil {
		fmt.Println(err)
		t.Errorf("error during open msg decoding: %v\n", err)
	}
	fmt.Printf("%#v\n", openMsg)
}

func TestEncodeMPcapability(t *testing.T) {
	mpCap := MPCapability{AFI: 1, SAFI: 1}
	encMpCap, err := EncodeMPCapability(mpCap)
	if err != nil {
		t.Errorf("cant encode mpCap")
	}
	encCap, err := EncodeCapability(Capability{Code: CAPABILITY_MP_EXTENSION}, encMpCap)
	if err != nil {
		t.Errorf("cant encode capability")
	}
	capability, data, err := DecodeCapability(encCap)
	if capability.Code != CAPABILITY_MP_EXTENSION {
		t.Errorf("error during capability decoding")
	}
	if err != nil {
		t.Errorf("can decode encoded capability")
	}
	decMpCap, err := DecodeMPCapability(data)
	if err != nil {
		t.Errorf("cant decode encoded mp capability")
	}
	if decMpCap.AFI != mpCap.AFI || decMpCap.SAFI != mpCap.SAFI {
		t.Errorf("error during enc/dec of mp cap")
	}
}

func TestEncodeOpenWithMPcapability(t *testing.T) {
	capList := []MPCapability{
		MPCapability{AFI: MP_AFI_IPV4, SAFI: MP_SAFI_UCAST},
		MPCapability{AFI: MP_AFI_IPV6, SAFI: MP_SAFI_UCAST}}
	openMsg := OpenMsg{Hdr: OpenMsgHdr{Version: 4, MyASN: 65000, HoldTime: 90, BGPID: 167772162}}
	openMsg.MPCaps = append(openMsg.MPCaps, capList...)
	data, err := EncodeOpenMsg(&openMsg)
	if err != nil {
		t.Errorf("cant encode open msg: %v\n", err)
	}
	_, err = DecodeOpenMsg(data[MSG_HDR_SIZE:])
	if err != nil {
		t.Errorf("cant decoded encoded msg: %v\n", err)
	}
}

func TestEncodeOpenMsg(t *testing.T) {
	encodedOpen, _ := hex.DecodeString(hexOpenMsg)
	openMsg := OpenMsg{Hdr: OpenMsgHdr{Version: 4, MyASN: 65000, HoldTime: 90, BGPID: 167772162, OptParamLength: 30}}
	encOpenMsg, err := EncodeOpenMsg(&openMsg)
	if err != nil {
		fmt.Println(err)
		t.Errorf("error during open msg  encoding")
	}
	//HACKISH TEST; we dont know how to encode all of the opt params and caps in etalon msg
	//so here we only tests how we have encoded ans,holdtime etc
	for cntr := 19; cntr < MIN_OPEN_MSG_SIZE-2; cntr++ {
		if encOpenMsg[cntr] != encodedOpen[cntr] {
			t.Errorf("byte of encoded msg is not equal to etalon's msg")
		}
	}
}

func TestDecodeUpdateMsg(t *testing.T) {
	encodedUpdate, _ := hex.DecodeString(hexUpdate1)
	_, err := DecodeUpdateMsg(encodedUpdate)
	if err != nil {
		fmt.Println(err)
		t.Errorf("error during update  msg decoding")
	}
	//PrintBgpUpdate(&bgpRoute)
	encodedUpdate, _ = hex.DecodeString(hexUpdate2)
	_, err = DecodeUpdateMsg(encodedUpdate)
	if err != nil {
		fmt.Println(err)
		t.Errorf("error during update  msg decoding")
	}
	//PrintBgpUpdate(&bgpRoute)
}

/*
   this test fails right now coz we dont support 32bit asn yet
   bogus in as_path part; should be 1 2 3 4 5 6
*/

func TestDecodeUpdMsgWithAsPath(t *testing.T) {
	encodedUpdate, _ := hex.DecodeString(hexUpdate4)
	_, err := DecodeUpdateMsg(encodedUpdate)
	if err != nil {
		fmt.Println(err)
		t.Errorf("error during update  msg decoding")
	}
	//PrintBgpUpdate(&bgpRoute)

}

func TestEncodeKeepaliveMsg(t *testing.T) {
	encodedKA, _ := hex.DecodeString(hexKA)
	encKA := GenerateKeepalive()
	for cntr := 0; cntr < len(encKA); cntr++ {
		if encKA[cntr] != encodedKA[cntr] {
			t.Errorf("byte of encoded msg is not equal to etalon's msg")
		}
	}
}

func TestDecodeNotificationMsg(t *testing.T) {
	encodedNotification, _ := hex.DecodeString(hexNotification)
	notification, err := DecodeNotificationMsg(encodedNotification)
	if err != nil {
		t.Errorf("error during notification decoding")
	}
	if notification.ErrorCode != 6 && notification.ErrorSubcode != 7 {
		t.Errorf("error during notification decoding(code and subcode are not equal to etalon)")
	}
}

func TestEncodeNotificationMsg(t *testing.T) {
	encodedNotification, _ := hex.DecodeString(hexNotification)
	notification := NotificationMsg{ErrorCode: BGP_CASE_ERROR, ErrorSubcode: BGP_CASE_ERROR_COLLISION}
	encNotification, err := EncodeNotificationMsg(&notification)
	if err != nil {
		fmt.Println(err)
		t.Errorf("error during notification encoding")
	}
	for cntr := 0; cntr < len(encNotification); cntr++ {
		if encNotification[cntr] != encodedNotification[cntr] {
			t.Errorf("byte of encoded msg is not equal to etalon's msg")
		}
	}

}

func TestEncodeUpdateMsg1(t *testing.T) {
	bgpRoute := BGPRoute{
		ORIGIN:          ORIGIN_IGP,
		MULTI_EXIT_DISC: uint32(123),
		LOCAL_PREF:      uint32(11),
		ATOMIC_AGGR:     true,
	}
	p1, _ := IPv4ToUint32("1.92.0.0")
	p2, _ := IPv4ToUint32("11.92.128.0")
	p3, _ := IPv4ToUint32("1.1.1.10")
	bgpRoute.Routes = append(bgpRoute.Routes, IPV4_NLRI{Length: 12, Prefix: p1})
	bgpRoute.Routes = append(bgpRoute.Routes, IPV4_NLRI{Length: 22, Prefix: p2})
	bgpRoute.Routes = append(bgpRoute.Routes, IPV4_NLRI{Length: 32, Prefix: p3})
	err := bgpRoute.AddV4NextHop("10.0.0.2")
	if err != nil {
		fmt.Println(err)
		t.Errorf("cant encode update msg")
	}
	data, err := EncodeUpdateMsg(&bgpRoute)
	if err != nil {
		fmt.Println(err)
		t.Errorf("cant encode update msg")
	}
	bgpRoute2, err := DecodeUpdateMsg(data)
	if err != nil {
		fmt.Println(err)
		t.Errorf("cant decode encoded update")
	}
	data2, _ := EncodeUpdateMsg(&bgpRoute2)
	if len(data) != len(data2) {
		t.Errorf("error in encoding/decoding of the same msg")
	}
	for cntr := 0; cntr < len(data); cntr++ {
		if data[cntr] != data2[cntr] {
			t.Errorf("error in encoding/decoding of the same msg")
			break
		}
	}
}

func TestEncodeWithdrawUpdateMsg1(t *testing.T) {
	bgpRoute := BGPRoute{}
	p4, _ := IPv4ToUint32("192.168.0.0")
	bgpRoute.WithdrawRoutes = append(bgpRoute.WithdrawRoutes, IPV4_NLRI{Length: 16, Prefix: p4})
	data, err := EncodeWithdrawUpdateMsg(&bgpRoute)
	if err != nil {
		fmt.Println(err)
		t.Errorf("cant encode withdraw update msg")
	}
	bgpRoute2, err := DecodeUpdateMsg(data)
	if err != nil {
		fmt.Println(err)
		t.Errorf("cant decode withdraw encoded update")
	}
	data2, _ := EncodeWithdrawUpdateMsg(&bgpRoute2)
	if len(data) != len(data2) {
		t.Errorf("error in encoding/decoding of the same withdraw msg")
	}
	for cntr := 0; cntr < len(data); cntr++ {
		if data[cntr] != data2[cntr] {
			t.Errorf("error in encoding/decoding of the same withdraw msg")
			break
		}
	}
}

func TestEncodeEndOfRIB(t *testing.T) {
	eor := GenerateEndOfRIB()
	if len(eor) != 23 {
		fmt.Println(eor)
		t.Errorf("error during EndOfRib marker generation")
	}
}

func TestIPv6StringToUint(t *testing.T) {
	_, err := IPv6StringToAddr("::")
	if err != nil {
		t.Errorf("cant convert ipv6 to ipv6addr\n")
	}
	addr, err := IPv6StringToAddr("fc1:2:3::1")
	if err != nil {
		t.Errorf("cant convert ipv6 to ipv6addr\n")
	}
	ipv6 := IPv6AddrToString(addr)
	fmt.Println(ipv6)
}

func TestIPv6NLRIEncoding(t *testing.T) {
	encodedIPv6NLRI, _ := hex.DecodeString(hexIPv6NLRI)
	nlri := IPV6_NLRI{Length: 48}
	v6addr, err := IPv6StringToAddr("2a00:bdc0:e003::")
	if err != nil {
		t.Errorf("error during ipv6 addr converting: %v\n", err)
	}
	nlri.Prefix = v6addr
	encIPv6NLRI, err := EncodeIPv6NLRI(nlri)
	if err != nil {
		t.Errorf("cant encode ipv6 nlri: %v\n", err)
	}
	fmt.Println(encodedIPv6NLRI)
	fmt.Println(encIPv6NLRI)
	if len(encodedIPv6NLRI) != len(encIPv6NLRI) {
		t.Errorf("len of encoded ipv6 nlri is not equal to len of etalon\n")
	}
	for i := 0; i < len(encIPv6NLRI); i++ {
		if encIPv6NLRI[i] != encodedIPv6NLRI[i] {
			t.Errorf("encoded ipv6 nlri is not equal to etalon")
		}
	}
}

func TestIPv6MP_REACH_Encoding(t *testing.T) {
	encodedIPv6MPREACH, _ := hex.DecodeString(hexIPv6_MP_REACH)
	nlri := IPV6_NLRI{Length: 48}
	v6addr, _ := IPv6StringToAddr("2a00:bdc0:e003::")
	v6nh, _ := IPv6StringToAddr("2001:7f8:20:101::245:180")
	nlri.Prefix = v6addr
	encIPv6MPREACH, err := EncodeIPV6_MP_REACH_NLRI(v6nh, nlri)
	if err != nil {
		t.Errorf("cant encode ipv6 mp reach nlri: %v\n", err)
	}
	if len(encodedIPv6MPREACH) != len(encIPv6MPREACH) {
		t.Errorf("len of encoded ipv6  mp reach nlri is not equal to len of etalon\n")
	}
	for i := 0; i < len(encIPv6MPREACH); i++ {
		if encIPv6MPREACH[i] != encodedIPv6MPREACH[i] {
			t.Errorf("encoded ipv6 mp reach nlri is not equal to etalon")
		}
	}
}

func TestIPv6MP_UNREACH_Encoding(t *testing.T) {
	nlri := IPV6_NLRI{Length: 48}
	v6addr, _ := IPv6StringToAddr("2a00:bdc0:e003::")
	nlri.Prefix = v6addr
	encIPv6MPUNREACH, err := EncodeIPV6_MP_UNREACH_NLRI(nlri)
	if err != nil {
		t.Errorf("cant encode ipv6 mp reach nlri: %v\n", err)
	}
	fmt.Println(encIPv6MPUNREACH)
}

func TestIPv6MP_REACH_PathAttrEncoding(t *testing.T) {
	encodedIPv6MPREACHPA, _ := hex.DecodeString(hexIPv6_MP_REACH_NLRI_PA)
	nlri := IPV6_NLRI{Length: 48}
	v6addr, _ := IPv6StringToAddr("2a00:bdc0:e003::")
	v6nh, _ := IPv6StringToAddr("2001:7f8:20:101::245:180")
	nlri.Prefix = v6addr
	pa := PathAttr{}
	encIPv6MPREACHPA, err := EncodeV6MPRNLRI(v6nh, nlri, &pa)
	if err != nil {
		t.Errorf("cant encode ipv6 mp reach nlri: %v\n", err)
	}
	if len(encodedIPv6MPREACHPA) != len(encIPv6MPREACHPA) {
		t.Errorf("len of encoded ipv6  mp reach nlri is not equal to len of etalon\n")
	}
	for i := 0; i < len(encIPv6MPREACHPA); i++ {
		if encIPv6MPREACHPA[i] != encodedIPv6MPREACHPA[i] {
			fmt.Println(encodedIPv6MPREACHPA)
			fmt.Println(encIPv6MPREACHPA)
			t.Errorf("encoded ipv6 mp reach nlri is not equal to etalon")
		}
	}
}

func TestIPv6MP_UNREACH_PathAttrEncoding(t *testing.T) {
	nlri := IPV6_NLRI{Length: 48}
	v6addr, _ := IPv6StringToAddr("2a00:bdc0:e003::")
	nlri.Prefix = v6addr
	pa := PathAttr{}
	encIPv6MPUNREACHPA, err := EncodeV6MPUNRNLRI(nlri, &pa)
	if err != nil {
		t.Errorf("cant encode ipv6 mp reach nlri: %v\n", err)
	}
	fmt.Println(encIPv6MPUNREACHPA)
}

func TestEncodeUpdateMsgV6(t *testing.T) {
	bgpRoute := BGPRoute{
		ORIGIN:          ORIGIN_IGP,
		MULTI_EXIT_DISC: uint32(123),
		LOCAL_PREF:      uint32(11),
		ATOMIC_AGGR:     true,
	}
	bgpRoute.NEXT_HOPv6, _ = IPv6StringToAddr("fc00::1")
	p1, _ := IPv6StringToAddr("2a02:6b8::")
	p2, _ := IPv6StringToAddr("2a00:1450:4010::")
	p3, _ := IPv6StringToAddr("2a03:2880:2130:cf05:face:b00c::1")
	bgpRoute.RoutesV6 = append(bgpRoute.RoutesV6, IPV6_NLRI{Length: 32, Prefix: p1})
	bgpRoute.RoutesV6 = append(bgpRoute.RoutesV6, IPV6_NLRI{Length: 48, Prefix: p2})
	bgpRoute.RoutesV6 = append(bgpRoute.RoutesV6, IPV6_NLRI{Length: 128, Prefix: p3})
	msg, err := EncodeUpdateMsg(&bgpRoute)
	if err != nil {
		t.Errorf("cant encode update msg with ipv6 mp_reach_nlri attr: %v\n", err)
	}
	fmt.Println(msg)
}

//Benchmarking

func BenchmarkDecodeUpdMsgWithAsPath(b *testing.B) {
	encodedUpdate, _ := hex.DecodeString(hexUpdate4)
	for i := 0; i < b.N; i++ {
		DecodeUpdateMsg(encodedUpdate)
	}
	//PrintBgpUpdate(&bgpRoute)

}

func BenchmarkEncodeUpdateMsg1(b *testing.B) {
	bgpRoute := BGPRoute{
		ORIGIN:          ORIGIN_IGP,
		MULTI_EXIT_DISC: uint32(123),
		LOCAL_PREF:      uint32(11),
		ATOMIC_AGGR:     true,
	}
	p1, _ := IPv4ToUint32("1.92.0.0")
	p2, _ := IPv4ToUint32("11.92.128.0")
	p3, _ := IPv4ToUint32("1.1.1.10")
	bgpRoute.Routes = append(bgpRoute.Routes, IPV4_NLRI{Length: 12, Prefix: p1})
	bgpRoute.Routes = append(bgpRoute.Routes, IPV4_NLRI{Length: 22, Prefix: p2})
	bgpRoute.Routes = append(bgpRoute.Routes, IPV4_NLRI{Length: 32, Prefix: p3})
	bgpRoute.AddV4NextHop("10.0.0.2")
	for i := 0; i < b.N; i++ {
		EncodeUpdateMsg(&bgpRoute)
	}
}

func BenchmarkEncodeUpdateMsgV6(b *testing.B) {
	bgpRoute := BGPRoute{
		ORIGIN:          ORIGIN_IGP,
		MULTI_EXIT_DISC: uint32(123),
		LOCAL_PREF:      uint32(11),
		ATOMIC_AGGR:     true,
	}
	bgpRoute.NEXT_HOPv6, _ = IPv6StringToAddr("fc00::1")
	p1, _ := IPv6StringToAddr("2a02:6b8::")
	p2, _ := IPv6StringToAddr("2a00:1450:4010::")
	p3, _ := IPv6StringToAddr("2a03:2880:2130:cf05:face:b00c::1")
	bgpRoute.RoutesV6 = append(bgpRoute.RoutesV6, IPV6_NLRI{Length: 32, Prefix: p1})
	bgpRoute.RoutesV6 = append(bgpRoute.RoutesV6, IPV6_NLRI{Length: 48, Prefix: p2})
	bgpRoute.RoutesV6 = append(bgpRoute.RoutesV6, IPV6_NLRI{Length: 128, Prefix: p3})
	for i := 0; i < b.N; i++ {
		EncodeUpdateMsg(&bgpRoute)
	}
}

func BenchmarkEncodeOpen(b *testing.B) {
	capList := []MPCapability{
		MPCapability{AFI: MP_AFI_IPV4, SAFI: MP_SAFI_UCAST},
		MPCapability{AFI: MP_AFI_IPV6, SAFI: MP_SAFI_UCAST}}
	openMsg := OpenMsg{Hdr: OpenMsgHdr{Version: 4, MyASN: 65000, HoldTime: 90, BGPID: 167772162}}
	openMsg.MPCaps = append(openMsg.MPCaps, capList...)
	for i := 0; i < b.N; i++ {
		EncodeOpenMsg(&openMsg)
	}
}

func BenchmarkDecodeOpen(b *testing.B) {
	capList := []MPCapability{
		MPCapability{AFI: MP_AFI_IPV4, SAFI: MP_SAFI_UCAST},
		MPCapability{AFI: MP_AFI_IPV6, SAFI: MP_SAFI_UCAST}}
	openMsg := OpenMsg{Hdr: OpenMsgHdr{Version: 4, MyASN: 65000, HoldTime: 90, BGPID: 167772162}}
	openMsg.MPCaps = append(openMsg.MPCaps, capList...)
	data, _ := EncodeOpenMsg(&openMsg)
	for i := 0; i < b.N; i++ {
		DecodeOpenMsg(data[MSG_HDR_SIZE:])
	}
}
