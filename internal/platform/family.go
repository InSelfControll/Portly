package platform

import "strings"

// determineFamily determines the OS family from ID and ID_LIKE fields
func determineFamily(id, idLike string) OSFamily {
	id = strings.ToLower(id)
	idLike = strings.ToLower(idLike)

	rhelVariants := []string{"rhel", "fedora", "centos", "rocky", "almalinux", "ol"}
	for _, v := range rhelVariants {
		if strings.Contains(id, v) || strings.Contains(idLike, v) {
			return FamilyRHEL
		}
	}

	debianVariants := []string{"debian", "ubuntu", "mint", "pop"}
	for _, v := range debianVariants {
		if strings.Contains(id, v) || strings.Contains(idLike, v) {
			return FamilyDebian
		}
	}

	return FamilyUnknown
}
