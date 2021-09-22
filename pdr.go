package gtp5gnl

import (
	"syscall"

	"github.com/khirono/go-genl"
	"github.com/khirono/go-nl"
)

// m LINK: ifindex
// m PDR_ID: u16
// o PDR_PRECEDENCE: u32
// o PDR_OUTER_HEADER_REMOVAL: u8
// o PDR_FAR_ID: u32
// o PDR_QER_ID: u32
// o PDR_PDI {
// o   PDI_UE_ADDR_IPV4: u32
// o   PDI_F_TEID {
//       F_TEID_I_TEID: u32
//       F_TEID_GTPU_ADDR_IPV4: u32
//     }
// o   PDI_SDF_FILTER {
//     }
//   }
func CreatePDR(c *Client, link *Link, pdrid int, attrs []nl.Attr) error {
	flags := syscall.NLM_F_EXCL
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_ADD_PDR})
	if err != nil {
		return err
	}
	err = req.Append(&nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  PDR_ID,
			Value: nl.AttrU16(pdrid),
		},
	})
	if err != nil {
		return err
	}
	err = req.Append(nl.AttrList(attrs))
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	return err
}

func UpdatePDR(c *Client, link *Link, pdrid int, attrs []nl.Attr) error {
	flags := syscall.NLM_F_REPLACE
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_ADD_PDR})
	if err != nil {
		return err
	}
	err = req.Append(&nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  PDR_ID,
			Value: nl.AttrU16(pdrid),
		},
	})
	if err != nil {
		return err
	}
	err = req.Append(nl.AttrList(attrs))
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	return err
}

func RemovePDR(c *Client, link *Link, pdrid int) error {
	flags := syscall.NLM_F_EXCL
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_DEL_PDR})
	if err != nil {
		return err
	}
	err = req.Append(&nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  PDR_ID,
			Value: nl.AttrU16(pdrid),
		},
	})
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	return err
}

func GetPDR(c *Client, link *Link, pdrid int) (*PDR, error) {
	flags := syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_PDR})
	if err != nil {
		return nil, err
	}
	err = req.Append(&nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  PDR_ID,
			Value: nl.AttrU16(pdrid),
		},
	})
	if err != nil {
		return nil, err
	}
	rsps, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if len(rsps) < 1 {
		return nil, err
	}
	pdr, err := DecodePDR(rsps[0].Body[genl.SizeofHeader:])
	if err != nil {
		return nil, err
	}
	return pdr, err
}

func GetPDRAll(c *Client) ([]PDR, error) {
	flags := syscall.NLM_F_DUMP
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_PDR})
	if err != nil {
		return nil, err
	}
	rsps, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	var pdrs []PDR
	for _, rsp := range rsps {
		pdr, err := DecodePDR(rsp.Body[genl.SizeofHeader:])
		if err != nil {
			return nil, err
		}
		pdrs = append(pdrs, *pdr)
	}
	return pdrs, err
}
