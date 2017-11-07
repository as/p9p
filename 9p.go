package p9p

type Kind byte

const (
	KindTversion, KindRversion Kind = iota + 100, iota + 101
	KindTauth, KindRauth
	KindTerror, KindRerror
	KindTflush, KindRflush
	KindTwalk, KindRwalk
	KindTopen, KindRopen
	KindTcreate, KindRcreate
	KindTread, KindRread
	KindTwrite, KindRwrite
	KindTclunk, KindRclunk
	KindTremove, KindRremove
	KindTstat, KindRstat
	KindTwstat, KindRwstat
)

//wire9 qid data[13]
//wire9 s n[2] data[n]
//wire9 Hdr size[4] msg[1]
//wire9 Tversion size[4] msg[1] tag[2] msize[4] version[,s]
//wire9 Rversion size[4] msg[1] tag[2] msize[4] version[,s]
//wire9 Tauth size[4] msg[1] tag[2] afid[4] uname[,s] aname[,s]
//wire9 Rauth size[4] msg[1] tag[2] aqid[13]
//wire9 Rerror size[4] msg[1] tag[2] ename[,s]
//wire9 Tflush size[4] msg[1] tag[2] oldtag[2]
//wire9 Rflush size[4] msg[1] tag[2]
//wire9 Tattach size[4] msg[1] tag[2] fid[4] afid[4] uname[,s] aname[,s]
//wire9 Rattach size[4] msg[1] tag[2] qid[13]
//wire9 Twalk size[4] msg[1] tag[2] fid[4] newfid[4] nwname[2] wname[nwname, []s]
//wire9 Rwalk size[4] msg[1] tag[2] nwqid[2] wqid[nwqid, []qid]
//wire9 Topen size[4] msg[1] tag[2] fid[4] mode[1]
//wire9 Ropen size[4] msg[1] tag[2] qid[13] iounit[4]
//wire9 Topenfd size[4] msg[1] tag[2] fid[4] mode[1]
//wire9 Ropenfd size[4] msg[1] tag[2] qid[13] iounit[4] unixfd[4]
//wire9 Tcreate size[4] msg[1] tag[2] fid[4] name[,s] perm[4] mode[1]
//wire9 Rcreate size[4] msg[1] tag[2] qid[13] iounit[4]
//wire9 Tread size[4] msg[1] tag[2] fid[4] offset[8] count[4]
//wire9 Rread size[4] msg[1] tag[2] count[4] data[count]
//wire9 Twrite size[4] msg[1] tag[2] fid[4] offset[8] count[4] data[count]
//wire9 Rwrite size[4] msg[1] tag[2] count[4]
//wire9 Tclunk size[4] msg[1] tag[2] fid[4]
//wire9 Rclunk size[4] msg[1] tag[2]
//wire9 Tremove size[4] msg[1] tag[2] fid[4]
//wire9 Rremove size[4] msg[1] tag[2]
//wire9 Tstat size[4] msg[1] tag[2] fid[4]
//wire9 Rstat size[4] msg[1] tag[2] stat[size-4-1-2]
//wire9 Twstat size[4] msg[1] tag[2] fid[4] stat[size-4-1-2]
//wire9 Rwstat size[4] msg[1] tag[2]
