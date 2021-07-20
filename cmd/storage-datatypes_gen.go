package cmd

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *DiskInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 13 {
		err = msgp.ArrayError{Wanted: 13, Got: zb0001}
		return
	}
	z.Total, err = dc.ReadUint64()
	if err != nil {
		err = msgp.WrapError(err, "Total")
		return
	}
	z.Free, err = dc.ReadUint64()
	if err != nil {
		err = msgp.WrapError(err, "Free")
		return
	}
	z.Used, err = dc.ReadUint64()
	if err != nil {
		err = msgp.WrapError(err, "Used")
		return
	}
	z.UsedInodes, err = dc.ReadUint64()
	if err != nil {
		err = msgp.WrapError(err, "UsedInodes")
		return
	}
	z.FreeInodes, err = dc.ReadUint64()
	if err != nil {
		err = msgp.WrapError(err, "FreeInodes")
		return
	}
	z.FSType, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "FSType")
		return
	}
	z.RootDisk, err = dc.ReadBool()
	if err != nil {
		err = msgp.WrapError(err, "RootDisk")
		return
	}
	z.Healing, err = dc.ReadBool()
	if err != nil {
		err = msgp.WrapError(err, "Healing")
		return
	}
	z.Endpoint, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Endpoint")
		return
	}
	z.MountPath, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "MountPath")
		return
	}
	z.ID, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	err = z.Metrics.DecodeMsg(dc)
	if err != nil {
		err = msgp.WrapError(err, "Metrics")
		return
	}
	z.Error, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Error")
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *DiskInfo) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 13
	err = en.Append(0x9d)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Total)
	if err != nil {
		err = msgp.WrapError(err, "Total")
		return
	}
	err = en.WriteUint64(z.Free)
	if err != nil {
		err = msgp.WrapError(err, "Free")
		return
	}
	err = en.WriteUint64(z.Used)
	if err != nil {
		err = msgp.WrapError(err, "Used")
		return
	}
	err = en.WriteUint64(z.UsedInodes)
	if err != nil {
		err = msgp.WrapError(err, "UsedInodes")
		return
	}
	err = en.WriteUint64(z.FreeInodes)
	if err != nil {
		err = msgp.WrapError(err, "FreeInodes")
		return
	}
	err = en.WriteString(z.FSType)
	if err != nil {
		err = msgp.WrapError(err, "FSType")
		return
	}
	err = en.WriteBool(z.RootDisk)
	if err != nil {
		err = msgp.WrapError(err, "RootDisk")
		return
	}
	err = en.WriteBool(z.Healing)
	if err != nil {
		err = msgp.WrapError(err, "Healing")
		return
	}
	err = en.WriteString(z.Endpoint)
	if err != nil {
		err = msgp.WrapError(err, "Endpoint")
		return
	}
	err = en.WriteString(z.MountPath)
	if err != nil {
		err = msgp.WrapError(err, "MountPath")
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	err = z.Metrics.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Metrics")
		return
	}
	err = en.WriteString(z.Error)
	if err != nil {
		err = msgp.WrapError(err, "Error")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *DiskInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 13
	o = append(o, 0x9d)
	o = msgp.AppendUint64(o, z.Total)
	o = msgp.AppendUint64(o, z.Free)
	o = msgp.AppendUint64(o, z.Used)
	o = msgp.AppendUint64(o, z.UsedInodes)
	o = msgp.AppendUint64(o, z.FreeInodes)
	o = msgp.AppendString(o, z.FSType)
	o = msgp.AppendBool(o, z.RootDisk)
	o = msgp.AppendBool(o, z.Healing)
	o = msgp.AppendString(o, z.Endpoint)
	o = msgp.AppendString(o, z.MountPath)
	o = msgp.AppendString(o, z.ID)
	o, err = z.Metrics.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Metrics")
		return
	}
	o = msgp.AppendString(o, z.Error)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DiskInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 13 {
		err = msgp.ArrayError{Wanted: 13, Got: zb0001}
		return
	}
	z.Total, bts, err = msgp.ReadUint64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Total")
		return
	}
	z.Free, bts, err = msgp.ReadUint64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Free")
		return
	}
	z.Used, bts, err = msgp.ReadUint64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Used")
		return
	}
	z.UsedInodes, bts, err = msgp.ReadUint64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "UsedInodes")
		return
	}
	z.FreeInodes, bts, err = msgp.ReadUint64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "FreeInodes")
		return
	}
	z.FSType, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "FSType")
		return
	}
	z.RootDisk, bts, err = msgp.ReadBoolBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "RootDisk")
		return
	}
	z.Healing, bts, err = msgp.ReadBoolBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Healing")
		return
	}
	z.Endpoint, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Endpoint")
		return
	}
	z.MountPath, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "MountPath")
		return
	}
	z.ID, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	bts, err = z.Metrics.UnmarshalMsg(bts)
	if err != nil {
		err = msgp.WrapError(err, "Metrics")
		return
	}
	z.Error, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Error")
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *DiskInfo) Msgsize() (s int) {
	s = 1 + msgp.Uint64Size + msgp.Uint64Size + msgp.Uint64Size + msgp.Uint64Size + msgp.Uint64Size + msgp.StringPrefixSize + len(z.FSType) + msgp.BoolSize + msgp.BoolSize + msgp.StringPrefixSize + len(z.Endpoint) + msgp.StringPrefixSize + len(z.MountPath) + msgp.StringPrefixSize + len(z.ID) + z.Metrics.Msgsize() + msgp.StringPrefixSize + len(z.Error)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *DiskMetrics) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "APILatencies":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "APILatencies")
				return
			}
			if z.APILatencies == nil {
				z.APILatencies = make(map[string]string, zb0002)
			} else if len(z.APILatencies) > 0 {
				for key := range z.APILatencies {
					delete(z.APILatencies, key)
				}
			}
			for zb0002 > 0 {
				zb0002--
				var za0001 string
				var za0002 string
				za0001, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "APILatencies")
					return
				}
				za0002, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "APILatencies", za0001)
					return
				}
				z.APILatencies[za0001] = za0002
			}
		case "APICalls":
			var zb0003 uint32
			zb0003, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "APICalls")
				return
			}
			if z.APICalls == nil {
				z.APICalls = make(map[string]uint64, zb0003)
			} else if len(z.APICalls) > 0 {
				for key := range z.APICalls {
					delete(z.APICalls, key)
				}
			}
			for zb0003 > 0 {
				zb0003--
				var za0003 string
				var za0004 uint64
				za0003, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "APICalls")
					return
				}
				za0004, err = dc.ReadUint64()
				if err != nil {
					err = msgp.WrapError(err, "APICalls", za0003)
					return
				}
				z.APICalls[za0003] = za0004
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *DiskMetrics) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "APILatencies"
	err = en.Append(0x82, 0xac, 0x41, 0x50, 0x49, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.APILatencies)))
	if err != nil {
		err = msgp.WrapError(err, "APILatencies")
		return
	}
	for za0001, za0002 := range z.APILatencies {
		err = en.WriteString(za0001)
		if err != nil {
			err = msgp.WrapError(err, "APILatencies")
			return
		}
		err = en.WriteString(za0002)
		if err != nil {
			err = msgp.WrapError(err, "APILatencies", za0001)
			return
		}
	}
	// write "APICalls"
	err = en.Append(0xa8, 0x41, 0x50, 0x49, 0x43, 0x61, 0x6c, 0x6c, 0x73)
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.APICalls)))
	if err != nil {
		err = msgp.WrapError(err, "APICalls")
		return
	}
	for za0003, za0004 := range z.APICalls {
		err = en.WriteString(za0003)
		if err != nil {
			err = msgp.WrapError(err, "APICalls")
			return
		}
		err = en.WriteUint64(za0004)
		if err != nil {
			err = msgp.WrapError(err, "APICalls", za0003)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *DiskMetrics) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "APILatencies"
	o = append(o, 0x82, 0xac, 0x41, 0x50, 0x49, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.APILatencies)))
	for za0001, za0002 := range z.APILatencies {
		o = msgp.AppendString(o, za0001)
		o = msgp.AppendString(o, za0002)
	}
	// string "APICalls"
	o = append(o, 0xa8, 0x41, 0x50, 0x49, 0x43, 0x61, 0x6c, 0x6c, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.APICalls)))
	for za0003, za0004 := range z.APICalls {
		o = msgp.AppendString(o, za0003)
		o = msgp.AppendUint64(o, za0004)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DiskMetrics) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "APILatencies":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "APILatencies")
				return
			}
			if z.APILatencies == nil {
				z.APILatencies = make(map[string]string, zb0002)
			} else if len(z.APILatencies) > 0 {
				for key := range z.APILatencies {
					delete(z.APILatencies, key)
				}
			}
			for zb0002 > 0 {
				var za0001 string
				var za0002 string
				zb0002--
				za0001, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "APILatencies")
					return
				}
				za0002, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "APILatencies", za0001)
					return
				}
				z.APILatencies[za0001] = za0002
			}
		case "APICalls":
			var zb0003 uint32
			zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "APICalls")
				return
			}
			if z.APICalls == nil {
				z.APICalls = make(map[string]uint64, zb0003)
			} else if len(z.APICalls) > 0 {
				for key := range z.APICalls {
					delete(z.APICalls, key)
				}
			}
			for zb0003 > 0 {
				var za0003 string
				var za0004 uint64
				zb0003--
				za0003, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "APICalls")
					return
				}
				za0004, bts, err = msgp.ReadUint64Bytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "APICalls", za0003)
					return
				}
				z.APICalls[za0003] = za0004
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *DiskMetrics) Msgsize() (s int) {
	s = 1 + 13 + msgp.MapHeaderSize
	if z.APILatencies != nil {
		for za0001, za0002 := range z.APILatencies {
			_ = za0002
			s += msgp.StringPrefixSize + len(za0001) + msgp.StringPrefixSize + len(za0002)
		}
	}
	s += 9 + msgp.MapHeaderSize
	if z.APICalls != nil {
		for za0003, za0004 := range z.APICalls {
			_ = za0004
			s += msgp.StringPrefixSize + len(za0003) + msgp.Uint64Size
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FileInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 25 {
		err = msgp.ArrayError{Wanted: 25, Got: zb0001}
		return
	}
	z.Volume, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Volume")
		return
	}
	z.Name, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	z.VersionID, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "VersionID")
		return
	}
	z.IsLatest, err = dc.ReadBool()
	if err != nil {
		err = msgp.WrapError(err, "IsLatest")
		return
	}
	z.Deleted, err = dc.ReadBool()
	if err != nil {
		err = msgp.WrapError(err, "Deleted")
		return
	}
	z.TransitionStatus, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "TransitionStatus")
		return
	}
	z.TransitionedObjName, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "TransitionedObjName")
		return
	}
	z.TransitionTier, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "TransitionTier")
		return
	}
	z.TransitionVersionID, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "TransitionVersionID")
		return
	}
	z.ExpireRestored, err = dc.ReadBool()
	if err != nil {
		err = msgp.WrapError(err, "ExpireRestored")
		return
	}
	z.DataDir, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "DataDir")
		return
	}
	z.XLV1, err = dc.ReadBool()
	if err != nil {
		err = msgp.WrapError(err, "XLV1")
		return
	}
	z.ModTime, err = dc.ReadTime()
	if err != nil {
		err = msgp.WrapError(err, "ModTime")
		return
	}
	z.Size, err = dc.ReadInt64()
	if err != nil {
		err = msgp.WrapError(err, "Size")
		return
	}
	z.Mode, err = dc.ReadUint32()
	if err != nil {
		err = msgp.WrapError(err, "Mode")
		return
	}
	var zb0002 uint32
	zb0002, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err, "Metadata")
		return
	}
	if z.Metadata == nil {
		z.Metadata = make(map[string]string, zb0002)
	} else if len(z.Metadata) > 0 {
		for key := range z.Metadata {
			delete(z.Metadata, key)
		}
	}
	for zb0002 > 0 {
		zb0002--
		var za0001 string
		var za0002 string
		za0001, err = dc.ReadString()
		if err != nil {
			err = msgp.WrapError(err, "Metadata")
			return
		}
		za0002, err = dc.ReadString()
		if err != nil {
			err = msgp.WrapError(err, "Metadata", za0001)
			return
		}
		z.Metadata[za0001] = za0002
	}
	var zb0003 uint32
	zb0003, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err, "Parts")
		return
	}
	if cap(z.Parts) >= int(zb0003) {
		z.Parts = (z.Parts)[:zb0003]
	} else {
		z.Parts = make([]ObjectPartInfo, zb0003)
	}
	for za0003 := range z.Parts {
		err = z.Parts[za0003].DecodeMsg(dc)
		if err != nil {
			err = msgp.WrapError(err, "Parts", za0003)
			return
		}
	}
	err = z.Erasure.DecodeMsg(dc)
	if err != nil {
		err = msgp.WrapError(err, "Erasure")
		return
	}
	z.MarkDeleted, err = dc.ReadBool()
	if err != nil {
		err = msgp.WrapError(err, "MarkDeleted")
		return
	}
	err = z.ReplicationState.DecodeMsg(dc)
	if err != nil {
		err = msgp.WrapError(err, "ReplicationState")
		return
	}
	z.Data, err = dc.ReadBytes(z.Data)
	if err != nil {
		err = msgp.WrapError(err, "Data")
		return
	}
	z.NumVersions, err = dc.ReadInt()
	if err != nil {
		err = msgp.WrapError(err, "NumVersions")
		return
	}
	z.SuccessorModTime, err = dc.ReadTime()
	if err != nil {
		err = msgp.WrapError(err, "SuccessorModTime")
		return
	}
	z.Fresh, err = dc.ReadBool()
	if err != nil {
		err = msgp.WrapError(err, "Fresh")
		return
	}
	z.Idx, err = dc.ReadInt()
	if err != nil {
		err = msgp.WrapError(err, "Idx")
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FileInfo) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 25
	err = en.Append(0xdc, 0x0, 0x19)
	if err != nil {
		return
	}
	err = en.WriteString(z.Volume)
	if err != nil {
		err = msgp.WrapError(err, "Volume")
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	err = en.WriteString(z.VersionID)
	if err != nil {
		err = msgp.WrapError(err, "VersionID")
		return
	}
	err = en.WriteBool(z.IsLatest)
	if err != nil {
		err = msgp.WrapError(err, "IsLatest")
		return
	}
	err = en.WriteBool(z.Deleted)
	if err != nil {
		err = msgp.WrapError(err, "Deleted")
		return
	}
	err = en.WriteString(z.TransitionStatus)
	if err != nil {
		err = msgp.WrapError(err, "TransitionStatus")
		return
	}
	err = en.WriteString(z.TransitionedObjName)
	if err != nil {
		err = msgp.WrapError(err, "TransitionedObjName")
		return
	}
	err = en.WriteString(z.TransitionTier)
	if err != nil {
		err = msgp.WrapError(err, "TransitionTier")
		return
	}
	err = en.WriteString(z.TransitionVersionID)
	if err != nil {
		err = msgp.WrapError(err, "TransitionVersionID")
		return
	}
	err = en.WriteBool(z.ExpireRestored)
	if err != nil {
		err = msgp.WrapError(err, "ExpireRestored")
		return
	}
	err = en.WriteString(z.DataDir)
	if err != nil {
		err = msgp.WrapError(err, "DataDir")
		return
	}
	err = en.WriteBool(z.XLV1)
	if err != nil {
		err = msgp.WrapError(err, "XLV1")
		return
	}
	err = en.WriteTime(z.ModTime)
	if err != nil {
		err = msgp.WrapError(err, "ModTime")
		return
	}
	err = en.WriteInt64(z.Size)
	if err != nil {
		err = msgp.WrapError(err, "Size")
		return
	}
	err = en.WriteUint32(z.Mode)
	if err != nil {
		err = msgp.WrapError(err, "Mode")
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Metadata)))
	if err != nil {
		err = msgp.WrapError(err, "Metadata")
		return
	}
	for za0001, za0002 := range z.Metadata {
		err = en.WriteString(za0001)
		if err != nil {
			err = msgp.WrapError(err, "Metadata")
			return
		}
		err = en.WriteString(za0002)
		if err != nil {
			err = msgp.WrapError(err, "Metadata", za0001)
			return
		}
	}
	err = en.WriteArrayHeader(uint32(len(z.Parts)))
	if err != nil {
		err = msgp.WrapError(err, "Parts")
		return
	}
	for za0003 := range z.Parts {
		err = z.Parts[za0003].EncodeMsg(en)
		if err != nil {
			err = msgp.WrapError(err, "Parts", za0003)
			return
		}
	}
	err = z.Erasure.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Erasure")
		return
	}
	err = en.WriteBool(z.MarkDeleted)
	if err != nil {
		err = msgp.WrapError(err, "MarkDeleted")
		return
	}
	err = z.ReplicationState.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "ReplicationState")
		return
	}
	err = en.WriteBytes(z.Data)
	if err != nil {
		err = msgp.WrapError(err, "Data")
		return
	}
	err = en.WriteInt(z.NumVersions)
	if err != nil {
		err = msgp.WrapError(err, "NumVersions")
		return
	}
	err = en.WriteTime(z.SuccessorModTime)
	if err != nil {
		err = msgp.WrapError(err, "SuccessorModTime")
		return
	}
	err = en.WriteBool(z.Fresh)
	if err != nil {
		err = msgp.WrapError(err, "Fresh")
		return
	}
	err = en.WriteInt(z.Idx)
	if err != nil {
		err = msgp.WrapError(err, "Idx")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FileInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 25
	o = append(o, 0xdc, 0x0, 0x19)
	o = msgp.AppendString(o, z.Volume)
	o = msgp.AppendString(o, z.Name)
	o = msgp.AppendString(o, z.VersionID)
	o = msgp.AppendBool(o, z.IsLatest)
	o = msgp.AppendBool(o, z.Deleted)
	o = msgp.AppendString(o, z.TransitionStatus)
	o = msgp.AppendString(o, z.TransitionedObjName)
	o = msgp.AppendString(o, z.TransitionTier)
	o = msgp.AppendString(o, z.TransitionVersionID)
	o = msgp.AppendBool(o, z.ExpireRestored)
	o = msgp.AppendString(o, z.DataDir)
	o = msgp.AppendBool(o, z.XLV1)
	o = msgp.AppendTime(o, z.ModTime)
	o = msgp.AppendInt64(o, z.Size)
	o = msgp.AppendUint32(o, z.Mode)
	o = msgp.AppendMapHeader(o, uint32(len(z.Metadata)))
	for za0001, za0002 := range z.Metadata {
		o = msgp.AppendString(o, za0001)
		o = msgp.AppendString(o, za0002)
	}
	o = msgp.AppendArrayHeader(o, uint32(len(z.Parts)))
	for za0003 := range z.Parts {
		o, err = z.Parts[za0003].MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "Parts", za0003)
			return
		}
	}
	o, err = z.Erasure.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Erasure")
		return
	}
	o = msgp.AppendBool(o, z.MarkDeleted)
	o, err = z.ReplicationState.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "ReplicationState")
		return
	}
	o = msgp.AppendBytes(o, z.Data)
	o = msgp.AppendInt(o, z.NumVersions)
	o = msgp.AppendTime(o, z.SuccessorModTime)
	o = msgp.AppendBool(o, z.Fresh)
	o = msgp.AppendInt(o, z.Idx)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FileInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 25 {
		err = msgp.ArrayError{Wanted: 25, Got: zb0001}
		return
	}
	z.Volume, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Volume")
		return
	}
	z.Name, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	z.VersionID, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "VersionID")
		return
	}
	z.IsLatest, bts, err = msgp.ReadBoolBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "IsLatest")
		return
	}
	z.Deleted, bts, err = msgp.ReadBoolBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Deleted")
		return
	}
	z.TransitionStatus, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "TransitionStatus")
		return
	}
	z.TransitionedObjName, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "TransitionedObjName")
		return
	}
	z.TransitionTier, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "TransitionTier")
		return
	}
	z.TransitionVersionID, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "TransitionVersionID")
		return
	}
	z.ExpireRestored, bts, err = msgp.ReadBoolBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "ExpireRestored")
		return
	}
	z.DataDir, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "DataDir")
		return
	}
	z.XLV1, bts, err = msgp.ReadBoolBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "XLV1")
		return
	}
	z.ModTime, bts, err = msgp.ReadTimeBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "ModTime")
		return
	}
	z.Size, bts, err = msgp.ReadInt64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Size")
		return
	}
	z.Mode, bts, err = msgp.ReadUint32Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Mode")
		return
	}
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Metadata")
		return
	}
	if z.Metadata == nil {
		z.Metadata = make(map[string]string, zb0002)
	} else if len(z.Metadata) > 0 {
		for key := range z.Metadata {
			delete(z.Metadata, key)
		}
	}
	for zb0002 > 0 {
		var za0001 string
		var za0002 string
		zb0002--
		za0001, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Metadata")
			return
		}
		za0002, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Metadata", za0001)
			return
		}
		z.Metadata[za0001] = za0002
	}
	var zb0003 uint32
	zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Parts")
		return
	}
	if cap(z.Parts) >= int(zb0003) {
		z.Parts = (z.Parts)[:zb0003]
	} else {
		z.Parts = make([]ObjectPartInfo, zb0003)
	}
	for za0003 := range z.Parts {
		bts, err = z.Parts[za0003].UnmarshalMsg(bts)
		if err != nil {
			err = msgp.WrapError(err, "Parts", za0003)
			return
		}
	}
	bts, err = z.Erasure.UnmarshalMsg(bts)
	if err != nil {
		err = msgp.WrapError(err, "Erasure")
		return
	}
	z.MarkDeleted, bts, err = msgp.ReadBoolBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "MarkDeleted")
		return
	}
	bts, err = z.ReplicationState.UnmarshalMsg(bts)
	if err != nil {
		err = msgp.WrapError(err, "ReplicationState")
		return
	}
	z.Data, bts, err = msgp.ReadBytesBytes(bts, z.Data)
	if err != nil {
		err = msgp.WrapError(err, "Data")
		return
	}
	z.NumVersions, bts, err = msgp.ReadIntBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "NumVersions")
		return
	}
	z.SuccessorModTime, bts, err = msgp.ReadTimeBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "SuccessorModTime")
		return
	}
	z.Fresh, bts, err = msgp.ReadBoolBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Fresh")
		return
	}
	z.Idx, bts, err = msgp.ReadIntBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Idx")
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FileInfo) Msgsize() (s int) {
	s = 3 + msgp.StringPrefixSize + len(z.Volume) + msgp.StringPrefixSize + len(z.Name) + msgp.StringPrefixSize + len(z.VersionID) + msgp.BoolSize + msgp.BoolSize + msgp.StringPrefixSize + len(z.TransitionStatus) + msgp.StringPrefixSize + len(z.TransitionedObjName) + msgp.StringPrefixSize + len(z.TransitionTier) + msgp.StringPrefixSize + len(z.TransitionVersionID) + msgp.BoolSize + msgp.StringPrefixSize + len(z.DataDir) + msgp.BoolSize + msgp.TimeSize + msgp.Int64Size + msgp.Uint32Size + msgp.MapHeaderSize
	if z.Metadata != nil {
		for za0001, za0002 := range z.Metadata {
			_ = za0002
			s += msgp.StringPrefixSize + len(za0001) + msgp.StringPrefixSize + len(za0002)
		}
	}
	s += msgp.ArrayHeaderSize
	for za0003 := range z.Parts {
		s += z.Parts[za0003].Msgsize()
	}
	s += z.Erasure.Msgsize() + msgp.BoolSize + z.ReplicationState.Msgsize() + msgp.BytesPrefixSize + len(z.Data) + msgp.IntSize + msgp.TimeSize + msgp.BoolSize + msgp.IntSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FileInfoVersions) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 4 {
		err = msgp.ArrayError{Wanted: 4, Got: zb0001}
		return
	}
	z.Volume, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Volume")
		return
	}
	z.Name, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	z.LatestModTime, err = dc.ReadTime()
	if err != nil {
		err = msgp.WrapError(err, "LatestModTime")
		return
	}
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err, "Versions")
		return
	}
	if cap(z.Versions) >= int(zb0002) {
		z.Versions = (z.Versions)[:zb0002]
	} else {
		z.Versions = make([]FileInfo, zb0002)
	}
	for za0001 := range z.Versions {
		err = z.Versions[za0001].DecodeMsg(dc)
		if err != nil {
			err = msgp.WrapError(err, "Versions", za0001)
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FileInfoVersions) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 4
	err = en.Append(0x94)
	if err != nil {
		return
	}
	err = en.WriteString(z.Volume)
	if err != nil {
		err = msgp.WrapError(err, "Volume")
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	err = en.WriteTime(z.LatestModTime)
	if err != nil {
		err = msgp.WrapError(err, "LatestModTime")
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Versions)))
	if err != nil {
		err = msgp.WrapError(err, "Versions")
		return
	}
	for za0001 := range z.Versions {
		err = z.Versions[za0001].EncodeMsg(en)
		if err != nil {
			err = msgp.WrapError(err, "Versions", za0001)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FileInfoVersions) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 4
	o = append(o, 0x94)
	o = msgp.AppendString(o, z.Volume)
	o = msgp.AppendString(o, z.Name)
	o = msgp.AppendTime(o, z.LatestModTime)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Versions)))
	for za0001 := range z.Versions {
		o, err = z.Versions[za0001].MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "Versions", za0001)
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FileInfoVersions) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 4 {
		err = msgp.ArrayError{Wanted: 4, Got: zb0001}
		return
	}
	z.Volume, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Volume")
		return
	}
	z.Name, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	z.LatestModTime, bts, err = msgp.ReadTimeBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "LatestModTime")
		return
	}
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Versions")
		return
	}
	if cap(z.Versions) >= int(zb0002) {
		z.Versions = (z.Versions)[:zb0002]
	} else {
		z.Versions = make([]FileInfo, zb0002)
	}
	for za0001 := range z.Versions {
		bts, err = z.Versions[za0001].UnmarshalMsg(bts)
		if err != nil {
			err = msgp.WrapError(err, "Versions", za0001)
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FileInfoVersions) Msgsize() (s int) {
	s = 1 + msgp.StringPrefixSize + len(z.Volume) + msgp.StringPrefixSize + len(z.Name) + msgp.TimeSize + msgp.ArrayHeaderSize
	for za0001 := range z.Versions {
		s += z.Versions[za0001].Msgsize()
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FilesInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Files":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "Files")
				return
			}
			if cap(z.Files) >= int(zb0002) {
				z.Files = (z.Files)[:zb0002]
			} else {
				z.Files = make([]FileInfo, zb0002)
			}
			for za0001 := range z.Files {
				err = z.Files[za0001].DecodeMsg(dc)
				if err != nil {
					err = msgp.WrapError(err, "Files", za0001)
					return
				}
			}
		case "IsTruncated":
			z.IsTruncated, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "IsTruncated")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *FilesInfo) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Files"
	err = en.Append(0x82, 0xa5, 0x46, 0x69, 0x6c, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Files)))
	if err != nil {
		err = msgp.WrapError(err, "Files")
		return
	}
	for za0001 := range z.Files {
		err = z.Files[za0001].EncodeMsg(en)
		if err != nil {
			err = msgp.WrapError(err, "Files", za0001)
			return
		}
	}
	// write "IsTruncated"
	err = en.Append(0xab, 0x49, 0x73, 0x54, 0x72, 0x75, 0x6e, 0x63, 0x61, 0x74, 0x65, 0x64)
	if err != nil {
		return
	}
	err = en.WriteBool(z.IsTruncated)
	if err != nil {
		err = msgp.WrapError(err, "IsTruncated")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *FilesInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Files"
	o = append(o, 0x82, 0xa5, 0x46, 0x69, 0x6c, 0x65, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Files)))
	for za0001 := range z.Files {
		o, err = z.Files[za0001].MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "Files", za0001)
			return
		}
	}
	// string "IsTruncated"
	o = append(o, 0xab, 0x49, 0x73, 0x54, 0x72, 0x75, 0x6e, 0x63, 0x61, 0x74, 0x65, 0x64)
	o = msgp.AppendBool(o, z.IsTruncated)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FilesInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Files":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Files")
				return
			}
			if cap(z.Files) >= int(zb0002) {
				z.Files = (z.Files)[:zb0002]
			} else {
				z.Files = make([]FileInfo, zb0002)
			}
			for za0001 := range z.Files {
				bts, err = z.Files[za0001].UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "Files", za0001)
					return
				}
			}
		case "IsTruncated":
			z.IsTruncated, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "IsTruncated")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *FilesInfo) Msgsize() (s int) {
	s = 1 + 6 + msgp.ArrayHeaderSize
	for za0001 := range z.Files {
		s += z.Files[za0001].Msgsize()
	}
	s += 12 + msgp.BoolSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *VolInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zb0001}
		return
	}
	z.Name, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	z.Created, err = dc.ReadTime()
	if err != nil {
		err = msgp.WrapError(err, "Created")
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z VolInfo) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 2
	err = en.Append(0x92)
	if err != nil {
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	err = en.WriteTime(z.Created)
	if err != nil {
		err = msgp.WrapError(err, "Created")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z VolInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	o = msgp.AppendString(o, z.Name)
	o = msgp.AppendTime(o, z.Created)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *VolInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zb0001}
		return
	}
	z.Name, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	z.Created, bts, err = msgp.ReadTimeBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Created")
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z VolInfo) Msgsize() (s int) {
	s = 1 + msgp.StringPrefixSize + len(z.Name) + msgp.TimeSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *VolsInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(VolsInfo, zb0002)
	}
	for zb0001 := range *z {
		var zb0003 uint32
		zb0003, err = dc.ReadArrayHeader()
		if err != nil {
			err = msgp.WrapError(err, zb0001)
			return
		}
		if zb0003 != 2 {
			err = msgp.ArrayError{Wanted: 2, Got: zb0003}
			return
		}
		(*z)[zb0001].Name, err = dc.ReadString()
		if err != nil {
			err = msgp.WrapError(err, zb0001, "Name")
			return
		}
		(*z)[zb0001].Created, err = dc.ReadTime()
		if err != nil {
			err = msgp.WrapError(err, zb0001, "Created")
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z VolsInfo) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0004 := range z {
		// array header, size 2
		err = en.Append(0x92)
		if err != nil {
			return
		}
		err = en.WriteString(z[zb0004].Name)
		if err != nil {
			err = msgp.WrapError(err, zb0004, "Name")
			return
		}
		err = en.WriteTime(z[zb0004].Created)
		if err != nil {
			err = msgp.WrapError(err, zb0004, "Created")
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z VolsInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zb0004 := range z {
		// array header, size 2
		o = append(o, 0x92)
		o = msgp.AppendString(o, z[zb0004].Name)
		o = msgp.AppendTime(o, z[zb0004].Created)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *VolsInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(VolsInfo, zb0002)
	}
	for zb0001 := range *z {
		var zb0003 uint32
		zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, zb0001)
			return
		}
		if zb0003 != 2 {
			err = msgp.ArrayError{Wanted: 2, Got: zb0003}
			return
		}
		(*z)[zb0001].Name, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, zb0001, "Name")
			return
		}
		(*z)[zb0001].Created, bts, err = msgp.ReadTimeBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, zb0001, "Created")
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z VolsInfo) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0004 := range z {
		s += 1 + msgp.StringPrefixSize + len(z[zb0004].Name) + msgp.TimeSize
	}
	return
}
