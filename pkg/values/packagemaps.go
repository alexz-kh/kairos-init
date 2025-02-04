package values

import (
	"bytes"
	"github.com/kairos-io/kairos-init/pkg/config"

	semver "github.com/hashicorp/go-version"
	sdkTypes "github.com/kairos-io/kairos-sdk/types"
)
import "text/template"

// packagemaps is a map of packages to install for each distro.
// so we can deal with stupid different names between distros.

// The format is usually a map[Distro]map[Architecture][]string
// So we can store the packages for each distro and architecture independently
// Except common packages, which are named the same across all distros
// Packages can be templated, so we can pass a map of parameters to replace in the package name
// So we can transform "linux-image-generic-hwe-{{.VERSION}}" into the proper version for each ubuntu release
// the params are not hardcoded or autogenerated anywhere yet.
// Ideally the System struct should have a method to generate the params for the packages automatically
// based on the distro and version, so we can pass them to the installer without anything from our side.
// Either we set also a Common key for the common packages, or we just duplicate them for both arches if needed
//

// CommonPackages are packages that are named the same across all distros and arches
var CommonPackages = []string{
	"file",       // Basic tool.
	"gawk",       // Basic tool.
	"iptables",   // Basic tool.
	"less",       // Basic tool.
	"nano",       // Basic tool.
	"sudo",       // Basic tool. Needed for the user to be able to run commands as root
	"tar",        // Basic tool.
	"zstd",       // Compression support for zstd
	"rsync",      // Install, upgrade, reset use it to sync the files
	"lvm2",       // Seems to be used to support rpi3 only
	"jq",         // No idea why we need it, check if we can drop it?
	"dosfstools", // For the fat32 partition on EFI systems
	"e2fsprogs",  // mkfs support for ext2/3/4
	"parted",     // Partitioning support, check if we need it anymore
}

// DistroFamilyInterface is an interface to get the value of a distro or family
// So we can refer to the package maps by the distro or family
type DistroFamilyInterface interface{}

type PackageMap map[DistroFamilyInterface]map[Architecture]VersionMap
type VersionMap map[string][]string

// ImmucorePackages are the minimum set of packages that immucore needs.
// Otherwise you wont be able to build the initrd with immucore on it.
var ImmucorePackages = PackageMap{
	DebianFamily: {
		ArchCommon: {
			Common: {
				"dracut",            // To build the initrd
				"dracut-network",    // Network-legacy support for dracut
				"isc-dhcp-common",   // Network-legacy support for dracut, basic tools
				"isc-dhcp-client",   // Network-legacy support for dracut, basic tools
				"cloud-guest-utils", // This brings growpart, so we can resize the partitions
			},
		},
	},
	Ubuntu: {
		ArchCommon: {
			">=22.04": {
				"dracut-live", // Livenet support for dracut, split into a separate package on 22.04
			},
		},
	},
	Debian: {
		ArchCommon: {
			Common: {
				"dracut-live",
			},
		},
	},
	RedHatFamily: {
		ArchCommon: {
			Common: {
				"dracut",
				"dracut-live",
				"dracut-network",
				"dracut-squash",
				"squashfs-tools",
				"dhcp-client",
			},
		},
	},
	SUSEFamily: {
		ArchCommon: {
			Common: {
				"dracut",
				"squashfs",
				"dhcp-client",
			},
		},
	},
}

// KernelPackages is a map of packages to install for each distro.
// No arch required here, maybe models will need different packages?
var KernelPackages = PackageMap{
	Ubuntu: {
		ArchCommon: {
			">=20.04, != 24.10": {
				// This is a template, so we can replace the version with the actual version of the system
				"linux-image-generic-hwe-{{.version}}",
			},
			// Somehow 24.10 uses the 22.04 hwe kernel
			"24.10": {"linux-image-generic-hwe-24.04"},
		},
	},
	Debian: {
		ArchAMD64: {
			Common: {
				"linux-image-amd64",
				"firmware-linux-free",
			},
		},
		ArchARM64: {
			Common: {
				"linux-image-arm64",
				"firmware-linux-free",
			},
		},
	},
	RedHatFamily: {
		ArchCommon: {
			Common: {
				"kernel",
				"kernel-modules",
				"kernel-modules-extra",
			},
		},
	},
	AlpineFamily: {
		ArchCommon: {
			Common: {
				"linux-lts",
			},
		},
	},
	SUSEFamily: {
		ArchCommon: {
			Common: {
				"kernel-default",
			},
		},
	},
}

// KernelPackagesTrustedBoot Separated kernel package for trusted boot as we dont want to install the same packages on both variants
// we need to keep teh trusted boot variant as small as possible so we want more control over it
// In this case, only Ubuntu has an specific smallest kernel package as its the only distro that supports trusted boot
// Fedora also works but we havent make it slim yet
var KernelPackagesTrustedBoot = PackageMap{
	Debian: {
		ArchAMD64: {
			Common: {
				"linux-image-amd64",
				"firmware-linux-free",
			},
		},
		ArchARM64: {
			Common: {
				"linux-image-arm64",
				"firmware-linux-free",
			},
		},
	},
	RedHatFamily: {
		ArchCommon: {
			Common: {
				"kernel",
				"kernel-modules",
				"kernel-modules-extra",
			},
		},
	},
	AlpineFamily: {
		ArchCommon: {
			Common: {
				"linux-lts",
			},
		},
	},
	SUSEFamily: {
		ArchCommon: {
			Common: {
				"kernel-default",
			},
		},
	},
}

// BasePackages is a map of packages to install for each distro and architecture.
// This comprises the base packages that are needed for the system to work on a Kairos system
var BasePackages = PackageMap{
	DebianFamily: {
		ArchCommon: {
			Common: {
				"ca-certificates", // Basic certificates for secure communication
				"curl",            // Basic tool. Also needed for netbooting as it is used to download the netboot artifacts. On rockylinux conflicts with curl-minimal
				"binutils",
				"conntrack",
				"console-setup",
				"coreutils",
				"cryptsetup",
				"debianutils",
				"ethtool",
				"fuse3",
				"gdisk",
				"gnupg",
				"gnupg1-l10n",
				"haveged",
				"iproute2",
				"iptables",
				"iputils-ping",
				"krb5-locales",
				"libatm1",
				"libglib2.0-data",
				"libgpm2",
				"libldap-common",
				"libnss-systemd",
				"libpam-cap",
				"libsasl2-modules",
				"mdadm",
				"nbd-client",
				"ncurses-term",
				"neovim",
				"nfs-common",
				"nftables",
				"open-iscsi",
				"openssh-server",
				"open-vm-tools",
				"os-prober",
				"patch",
				"pigz",
				"pkg-config",
				"psmisc",
				"publicsuffix",
				"python3-pynvim",
				"shared-mime-info",
				"snapd",
				"systemd", // Basic tool.
				"systemd-timesyncd",
				"systemd-sysv", // provides reboot and shutdown commands. like what? they are just symlinks to systemctl lol
				"xauth",
				"xclip",
				"xdg-user-dirs",
				"xxd",
				"xz-utils",
				"zerofree",
			},
		},
	},
	SUSEFamily: {
		ArchCommon: {
			Common: {
				"curl", // Basic tool. Also needed for netbooting as it is used to download the netboot artifacts. On rockylinux conflicts with curl-minimal
				"bash-completion",
				"conntrack-tools",
				"cryptsetup",
				"coreutils",
				"device-mapper",
				"fail2ban",
				"findutils",
				"growpart",
				"gptfdisk",
				"haveged",
				"htop",
				"iproute2",
				"iputils",
				"issue-generator",
				"logrotate",
				"lsscsi",
				"mdadm",
				"multipath-tools",
				// "nfs-utils", // Not available by default, coming from extra repo.
				// "nohang", // Not available by default, coming from extra repo.
				"open-iscsi",
				"openssh",
				"open-vm-tools",
				"pigz",
				"policycoreutils",
				"polkit",
				"procps",
				"qemu-guest-agent", // TODO: Move this to generic model?
				"strace",
				"systemd",
				"systemd-network",
				"timezone",
				"tmux",
				"vim",
				"which",
				"tpm2*",
			},
		},
	},
	AlpineFamily: {
		ArchCommon: {
			Common: {
				"curl", // Basic tool. Also needed for netbooting as it is used to download the netboot artifacts. On rockylinux conflicts with curl-minimal
				"bash",
				"bash-completion",
				"blkid",
				"cloud-utils-growpart",
				"bonding",
				"bridge",
				"busybox-openrc",
				"ca-certificates",
				"connman",
				"conntrack-tools",
				"coreutils",
				"cryptsetup",
				"device-mapper-udev",
				"dbus",
				"dmidecode",
				"dosfstools",
				"e2fsprogs",
				"e2fsprogs-extra",
				"efibootmgr",
				"eudev",
				"eudev-hwids",
				"fail2ban",
				"findutils",
				"findmnt",
				"gcompat",
				"gettext",
				"haveged",
				"htop",
				"hvtools",
				"iproute2",
				"irqbalance",
				"iscsi-scst",
				"kbd-bkeymaps",
				"libc6-compat",
				"libusb",
				"lm-sensors",
				"logrotate",
				"lsscsi",
				"lvm2-extra",
				"mdadm",
				"mdadm-misc",
				"mdadm-udev",
				"multipath-tools",
				"ncurses",
				"ncurses-terminfo",
				"nfs-utils",
				"open-iscsi",
				"openrc",
				"openssh-client",
				"openssh-server",
				"open-vm-tools",
				"open-vm-tools-deploypkg",
				"open-vm-tools-guestinfo",
				"open-vm-tools-static",
				"open-vm-tools-vmbackup",
				"procps",
				"qemu-guest-agent",
				"rbd-nbd",
				"sgdisk",
				"smartmontools",
				"squashfs-tools",
				"strace",
				"tzdata",
				"util-linux",
				"vim",
				"which",
				"wireguard-tools",
				"wpa_supplicant",
				"xfsprogs",
				"xfsprogs-extra",
				"xz",
			},
		},
	},
	RedHatFamily: {
		ArchCommon: {
			Common: {
				"gdisk",                // Yip requires it for partitioning, maybe BasePackages
				"audit",                // For audit support, check if needed?
				"cracklib-dicts",       // Password dictionary support
				"cloud-utils-growpart", // grow partition use. Check if yip still needs it?
				"device-mapper",        // Device mapper support, needed for lvm and cryptsetup
				"openssh-server",
				"openssh-clients",
				"polkit",
				"qemu-guest-agent",
				"systemd", // Basic tool.
				"systemd-resolved",
				"which",      // Basic tool. Basepackages?
				"cryptsetup", // For encrypted partitions support, needed for trusted boot and dracut building
			},
		},
	},
	Debian: {
		ArchCommon: {
			Common: {
				"systemd-resolved",
				"nohang",
				"polkitd",
			},
			">=13": {
				"systemd-cryptsetup", // separated package on testing, so we need to add it on 13 and above
			},
		},
	},
	Ubuntu: {
		ArchCommon: {
			Common: {
				// TODO: Check if we need all of these packages, some of them are probably not needed or can go into the family?
				"fdisk", // Yip requires it for partitioning
				"conntrack",
				"console-data",      // Console font support
				"cloud-guest-utils", // Yip requires it, this brings growpart, so we can resize the partitions
				"gettext",
				"systemd-container",      // Not sure if needed?
				"ubuntu-advantage-tools", // For ubuntu advantage support, enablement of ubuntu services
				"tpm2-tools",             // For TPM support, mainly trusted boot
				"dmsetup",                // Device mapper support, needed for lvm and cryptsetup
				"networkd-dispatcher",
				"packagekit-tools",
				"publicsuffix",
				"xdg-user-dirs",
				"zfsutils-linux", // For zfs tools (zfs and zpool)
			},
			">=24.04": {
				"systemd-resolved", // For systemd-resolved support, added as a separate package on 24.04
			},
		},
	},
	Fedora: {
		ArchCommon: {
			Common: {
				"haveged",          // Random number generator, check if needed?
				"systemd-networkd", // Not available in other distros, too old version maybe?
			},
		},
	},
}

// GrubPackages is a map of packages to install for each distro and architecture.
// TODO: Check why some packages we only install on amd64 and not on arm64?? Like neovim???
// Note: some of the packages seems to be only installed here as we dont have any size restraints
// And we dont want to have Trusted Boot have those packages, as we want it small.
// we should probably move those into a new PackageMap called ExtendedPackages or something like that
// instead of merging them with grub packages.
var GrubPackages = PackageMap{
	DebianFamily: {
		ArchCommon: {
			Common: {
				"kbd",            // Keyboard configuration
				"lldpd",          // For lldp support, check if needed?
				"shim-signed",    // For secure boot support
				"snmpd",          // For snmp support, check if needed? Move to BasePackages if so?
				"squashfs-tools", // For squashfs support, probably needs to be part of BasePackages
				//"zfsutils-linux",        // For zfs tools (zfs and zpool), probably needs to be part of BasePackages
				// Requires a repo add for Debian
			},
		},
		ArchAMD64: {
			Common: {
				"grub2",                 // Basic grub support
				"grub-efi-amd64-bin",    // Basic grub support for EFI
				"grub-efi-amd64-signed", // For secure boot support
				"grub-pc-bin",           // Basic grub support for BIOS, probably needed byt AuroraBoot to build hybrid isos?
				"grub2-common",          // Basic grub support
			},
		},
		ArchARM64: {
			Common: {
				"grub-efi-arm64",        // Basic grub support for EFI
				"grub-efi-arm64-bin",    // Basic grub support for EFI
				"grub-efi-arm64-signed", // For secure boot support
			},
		},
	},
	RedHatFamily: {
		ArchCommon: {
			Common: {
				"grub2",
			},
		},
		ArchAMD64: {
			Common: {
				"grub2-efi-x64",
				"grub2-efi-x64-modules",
				"grub2-pc",
				"shim-x64",
			},
		},
		ArchARM64: {
			Common: {
				"grub2-efi-aa64",
				"grub2-efi-aa64-modules",
				"shim-aa64",
			},
		},
	},
	AlpineFamily: {
		ArchAMD64: {
			Common: {
				"grub-bios",
			},
		},
		ArchCommon: {
			Common: {
				"grub",
				"grub-efi",
			},
		},
	},
	SUSEFamily: {
		ArchCommon: {
			Common: {
				"nethogs",
				"patch",
				"shim",
				"iw",
			},
		},
		ArchAMD64: {
			Common: {
				"grub2-i386-pc",
				"grub2-x86_64-efi",
				"kernel-firmware-all",
			},
		},
		ArchARM64: {
			Common: {
				"bcm43xx-firmware",
				"grub2-arm64-efi",
				"kernel-firmware-ath10k",
				"kernel-firmware-ath11k",
				"kernel-firmware-atheros",
				"kernel-firmware-bluetooth",
				"kernel-firmware-brcm",
				"kernel-firmware-iwlwifi",
				"kernel-firmware-network",
				"kernel-firmware-realtek",
				"kernel-firmware-serial",
				"kernel-firmware-usb-network",
			},
		},
	},
	Ubuntu: {
		ArchCommon: {
			Common: {
				"zfsutils-linux", // For zfs tools (zfs and zpool) We put it here because this package is useless for the Trusted Boot variant and we want to keep that part as ligth as possible
			},
		},
	},
}

// SystemdPackages is a map of packages to install for each distro and architecture for systemd-boot (trusted boot) variants
// TODO: Check why some packages we only install on amd64 and not on arm64?? Like kmod???
var SystemdPackages = PackageMap{
	Ubuntu: {
		ArchCommon: {
			Common: {
				"systemd",
			},
			">=24.04": {
				"iucode-tool",
				"kmod",
				"linux-base",
				"systemd-boot", // Trusted boot support, it was split as a package on 24.04
			},
		},
	},
}

// RpiPackages is a map of packages to install for each distro and architecture for Raspberry Pi variants
// TODO: Actually implement this somehow somewhere lol
// TODO: Make it a board thing not only rpi
// TODO(debian): Needs to run `sed -i 's/^Components: main.*$/& non-free-firmware/' /etc/apt/sources.list.d/debian.sources` before installing the firmware
var RpiPackages = PackageMap{
	Debian: {
		ArchAMD64: {
			Rpi4.String(): {
				"raspi-firmware",
			},
		},
	},
	Arch: {
		ArchARM64: {
			Rpi3.String(): {
				"linux-rpi",
			},
			Rpi4.String(): {
				"linux-rpi4",
			},
		},
	},
	SUSEFamily: {
		ArchARM64: {
			Rpi3.String(): {
				"raspberrypi-eeprom",
				"raspberrypi-firmware",
				"raspberrypi-firmware-dt",
				"sysconfig",
				"sysconfig-netconfig",
				"sysvinit-tools",
				"wireless-tools",
				"wpa_supplicant",
			},
			Rpi4.String(): {
				"raspberrypi-eeprom",
				"raspberrypi-firmware",
				"raspberrypi-firmware-dt",
				"sysconfig",
				"sysconfig-netconfig",
				"sysvinit-tools",
				"wireless-tools",
				"wpa_supplicant",
			},
		},
	},
}

// PackageListToTemplate takes a list of packages and a map of parameters to replace in the package name
// and returns a list of packages with the parameters replaced.
func PackageListToTemplate(packages []string, params map[string]string, l sdkTypes.KairosLogger) ([]string, error) {
	var finalPackages []string
	for _, pkg := range packages {
		var result bytes.Buffer
		tmpl, err := template.New("versionTemplate").Parse(pkg)
		if err != nil {
			l.Logger.Error().Err(err).Str("package", pkg).Msg("Error parsing template.")
			return []string{}, err
		}
		err = tmpl.Execute(&result, params)
		if err != nil {
			l.Logger.Error().Err(err).Str("package", pkg).Msg("Error executing template.")
			return []string{}, err
		}
		finalPackages = append(finalPackages, result.String())
	}
	return finalPackages, nil
}

func GetPackages(s System, l sdkTypes.KairosLogger) ([]string, error) {
	mergedPkgs := CommonPackages

	// Go over all packages maps
	filteredPackages := []VersionMap{
		BasePackages[s.Distro][ArchCommon], // Common packages to both arches
		BasePackages[s.Family][ArchCommon], // Common packages to both arches by family
		BasePackages[s.Distro][s.Arch],     // Specific packages for the arch
		BasePackages[s.Family][s.Arch],     // Specific packages for the arch by family
	}

	if config.DefaultConfig.TrustedBoot {
		filteredPackages = append(filteredPackages, KernelPackagesTrustedBoot[s.Distro][ArchCommon]) // Common kernel packages to both arches
		filteredPackages = append(filteredPackages, KernelPackagesTrustedBoot[s.Family][ArchCommon]) // Common kernel packages to both arches by family
		filteredPackages = append(filteredPackages, KernelPackagesTrustedBoot[s.Distro][s.Arch])     // Specific kernel packages for the arch
		filteredPackages = append(filteredPackages, KernelPackagesTrustedBoot[s.Family][s.Arch])     // Specific kernel packages for the arch by family
		// Install only systemd-boot packages
		filteredPackages = append(filteredPackages, SystemdPackages[s.Distro][ArchCommon])
		filteredPackages = append(filteredPackages, SystemdPackages[s.Family][ArchCommon])
		filteredPackages = append(filteredPackages, SystemdPackages[s.Distro][s.Arch])
		filteredPackages = append(filteredPackages, SystemdPackages[s.Family][s.Arch])
	} else {
		filteredPackages = append(filteredPackages, KernelPackages[s.Distro][ArchCommon]) // Common kernel packages to both arches
		filteredPackages = append(filteredPackages, KernelPackages[s.Family][ArchCommon]) // Common kernel packages to both arches by family
		filteredPackages = append(filteredPackages, KernelPackages[s.Distro][s.Arch])     // Specific kernel packages for the arch
		filteredPackages = append(filteredPackages, KernelPackages[s.Family][s.Arch])     // Specific kernel packages for the arch by family
		// install grub and immucore packages
		filteredPackages = append(filteredPackages, GrubPackages[s.Distro][ArchCommon])
		filteredPackages = append(filteredPackages, GrubPackages[s.Family][ArchCommon])
		filteredPackages = append(filteredPackages, GrubPackages[s.Distro][s.Arch])
		filteredPackages = append(filteredPackages, GrubPackages[s.Family][s.Arch])
		filteredPackages = append(filteredPackages, ImmucorePackages[s.Distro][ArchCommon])
		filteredPackages = append(filteredPackages, ImmucorePackages[s.Family][ArchCommon])
		filteredPackages = append(filteredPackages, ImmucorePackages[s.Distro][s.Arch])
		filteredPackages = append(filteredPackages, ImmucorePackages[s.Family][s.Arch])
	}

	mergedPkgs = append(mergedPkgs, FilterPackagesOnConstraint(s, l, filteredPackages)...)

	return mergedPkgs, nil
}

// FilterPackagesOnConstraint filters the packages based on the system version and the constraints in the package map
func FilterPackagesOnConstraint(s System, l sdkTypes.KairosLogger, pkgsToFilter []VersionMap) []string {
	// Go over each list of packages
	var pkgs []string
	systemVersion, err := semver.NewVersion(s.Version)
	if err != nil {
		return pkgs
	}
	for _, packages := range pkgsToFilter {
		// for each package map, check if the version matches the constraint
		for constraint, values := range packages {
			// Add them if they are common
			l.Logger.Debug().Str("constraint", constraint).Str("version", systemVersion.String()).Msg("Checking constraint")
			if constraint == Common {
				l.Logger.Debug().Strs("packages", values).Msg("Adding common packages")
				pkgs = append(pkgs, values...)
				continue
			}
			semverConstraint, err := semver.NewConstraint(constraint)
			if err != nil {
				l.Logger.Error().Err(err).Str("constraint", constraint).Msg("Error parsing constraint.")
				continue
			}
			// Also add them if the constraint matches
			if semverConstraint.Check(systemVersion) {
				l.Logger.Debug().Strs("packages", values).Msg("Constraint matches, adding packages")
				pkgs = append(pkgs, values...)
			}
		}
	}
	return pkgs
}
