package examples

//type Ip struct {
//	Sip   uint32
//	Dip   uint32
//	SPort uint16
//	DPort uint16
//}
//
//func initAllowIpMap(m *ebpf.Map) {
//	ip := binary.BigEndian.Uint32(net.ParseIP("172.20.1.2").To4())
//	err := m.Put(ip, uint8(1))
//	if err != nil {
//		log.Fatalln(err)
//	}
//}
//
//func Load() {
//	xdpObj := xdpObjects{}
//
//	err := loadXdpObjects(&xdpObj, nil)
//	if err != nil {
//
//		log.Fatalln(err)
//	}
//
//	defer xdpObj.Close()
//
//	iface, err := net.InterfaceByName("docker0")
//	if err != nil {
//
//		log.Fatalln(err)
//	}
//	l, err := link.AttachXDP(link.XDPOptions{
//		Program:   xdpObj.MyPass,
//		Interface: iface.Index,
//	})
//	if err != nil {
//		log.Fatalln(err)
//	}
//	defer l.Close()
//
//	initAllowIpMap(xdpObj.IpMap)
//	reader, err := ringbuf.NewReader(xdpObj.IpMap)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	defer reader.Close()
//
//	go func() {
//		for {
//			record, err := reader.Read()
//			if err != nil {
//				if errors.Is(err, perf.ErrClosed) {
//					log.Println("Receiver signal, exiting...")
//					return
//				}
//				log.Println("reading from reader:", err)
//				continue
//			}
//			pointer := (*Ip)(unsafe.Pointer(&record.RawSample[0]))
//			sip := net.IPv4(byte(pointer.Sip), byte(pointer.Sip>>8), byte(pointer.Sip>>16), byte(pointer.Sip>>24))
//			dip := net.IPv4(byte(pointer.Dip), byte(pointer.Dip>>8), byte(pointer.Dip>>16), byte(pointer.Dip>>24))
//			log.Println("Record:", sip, dip)
//		}
//	}()
//
//	time.Sleep(time.Minute)
//}
