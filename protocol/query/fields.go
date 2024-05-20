package query

type Fields []Field

type Field string

const (
	FName    Field = "name"
	FExists  Field = "exists"
	FCclock  Field = "cclock"
	FOclock  Field = "oclock"
	FCtime   Field = "ctime"
	FCtimeMs Field = "ctime_ms"
	FCtimeUs Field = "ctime_us"
	FCtimeNs Field = "ctime_ns"
	FCtimeF  Field = "ctime_f"
	FMtime   Field = "mtime"
	FMtimeMs Field = "mtime_ms"
	FMtimeUs Field = "mtime_us"
	FMtimeNs Field = "mtime_ns"
	FMtimeF  Field = "mtime_f"
	FSize    Field = "size"
	FMode    Field = "mode"
	FUid     Field = "uid"
	FGid     Field = "gid"
	FIno     Field = "ino"
	FDev     Field = "dev"
	FNlink   Field = "nlink"
	FNew     Field = "new"
	FType    Field = "type"

	FSymlinkTarget  Field = "symlink_target"
	FContentSha1hex Field = "content.sha1hex"
)
