package model

// structure in https://gitlab.com/cryptsetup/LUKS2-docs.
type CryptsetupMeta struct {
	KeySlots map[string]struct {
		Type                 string `json:"type"`
		KeySize              int    `json:"key_size"`
		AntiForensicSplitter struct {
			Type    string `json:"type"`
			Stripes int    `json:"stripes"`
			Hash    string `json:"hash"`
		} `json:"af"`
		Area struct {
			Type       string `json:"type"`
			Offset     string `json:"offset"`
			Size       string `json:"size"`
			Encryption string `json:"encryption"`
			KeySize    int    `json:"key_size"`
		} `json:"area"`
		KDF struct {
			Type   string `json:"type"`
			Time   int    `json:"time"`
			Memory int    `json:"memory"`
			CPUs   int    `json:"cpus"`
			Salt   string `json:"salt"`
		} `json:"kdf"`
	} `json:"keyslots"`
	Tokens   map[string]struct{} `json:"tokens"`
	Segments map[string]struct {
		Type       string   `json:"type"`
		Offset     string   `json:"offset"`
		Size       string   `json:"size"`
		Flags      []string `json:"flags,omitempty"`
		IVTweak    string   `json:"iv_tweak"`
		Encryption string   `json:"encryption"`
		SectorSize int      `json:"sector_size"`
		Integrity  struct {
			Type              string `json:"type"`
			JournalEncryption string `json:"journal_encryption"`
			JournalIntegrity  string `json:"journal_integrity"`
			KeySize           int    `json:"key_size"`
		} `json:"integrity,omitempty"`
	} `json:"segments"`
	Digests map[string]struct {
		Type       string   `json:"type"`
		Keyslots   []string `json:"keyslots"`
		Segments   []string `json:"segments"`
		Hash       string   `json:"hash"`
		Iterations int      `json:"iterations"`
		Salt       string   `json:"salt"`
		Digest     string   `json:"digest"`
	} `json:"digests"`
	Config struct {
		JSONSize     string `json:"json_size"`
		KeyslotsSize string `json:"keyslots_size"`
	}
}
