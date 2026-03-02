package enum

type PermissionsKind string

const (
	PermissionsKindUnspecified PermissionsKind = "PERMISSIONS_KIND_UNSPECIFIED"
	PermissionsKindLayout      PermissionsKind = "PERMISSIONS_KIND_LAYOUT"
	PermissionsKindNote        PermissionsKind = "PERMISSIONS_KIND_NOTE"
)

func (pk PermissionsKind) String() string {
	return string(pk)
}

func PermissionsKindFromString(s string) PermissionsKind {
	switch s {
	case PermissionsKindLayout.String():
		return PermissionsKindLayout
	case PermissionsKindNote.String():
		return PermissionsKindNote
	default:
		return PermissionsKindUnspecified
	}
}
