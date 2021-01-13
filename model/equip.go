package model

import "gorm.io/gorm"

type Equip struct {
	Serial               int64  `gorm:"primary_key;column:serial"`
	Quality              uint8  `gorm:"default:0;column:quality"`
	MinUseLevel          uint8  `gorm:"default:0;column:min_use_level"`
	MaxUseLevel          uint8  `gorm:"default:0;column:max_use_level"`
	MinDmgIn             uint16 `gorm:"default:0;column:min_dmg_in"`
	MinDmg               int16  `gorm:"default:0;column:min_dmg"`
	MaxDmg               int16  `gorm:"default:0;column:max_dmg"`
	Armor                int16  `gorm:"default:0;column:armor"`
	PotVal               int32  `gorm:"default:0;column:potval"`
	PotValUsed           int32  `gorm:"default:0;column:potval_used"`
	RoleAttrEffect       []byte `gorm:"column:role_att_effect"`
	XerInID              uint32 `gorm:"default:4294967295;column:xer_in_id"`
	XerOutID             uint32 `gorm:"default:4294967295;column:xer_out_id"`
	SpecelAttr           int8   `gorm:"default:0;column:special_att"`
	Appearance           int16  `gorm:"default:0;column:appearance"`
	Rein                 uint8  `gorm:"default:0;column:rein"`
	Savvy                uint8  `gorm:"default:0;column:savvy"`
	Fortune              uint8  `gorm:"default:0;column:fortune"`
	ColorID              int8   `gorm:"default:0;column:color_id"`
	AttALimMod           int16  `gorm:"default:0;column:att_a_lim_mod"`
	AttALimModPct        int16  `gorm:"default:0;column:att_a_lim_mod_pct"`
	PosyTimes            int8   `gorm:"default:0;column:posy_times"`
	PosyEffect           []byte `gorm:"column:posy_effect"`
	EngraveTimes         int8   `gorm:"default:0;column:engrave_times"`
	EngraveAtt           []byte `gorm:"column:engrave_att"`
	HoleNum              int8   `gorm:"default:0;column:hole_num"`
	HoleGemID            []byte `gorm:"column:hole_gem_id"`
	BrandLevl            int8   `gorm:"default:0;column:brand_level"`
	DerateVal            []byte `gorm:"column:derate_val"`
	XFulLevl             int8   `gorm:"default:0;column:x_ful_level"`
	HoleGemNess          []byte `gorm:"column:hole_gem_ness"`
	CanCut               int8   `gorm:"default:1;column:can_cut"`
	FlareVal             uint8  `gorm:"default:0;column:flare_val"`
	QualityModPct        int16  `gorm:"default:0;column:quality_mod_pct"`
	QualityModPctEx      uint16 `gorm:"default:0;column:quality_mod_pct_ex"`
	PotValModPct         uint16 `gorm:"default:10000;column:pot_val_mod_pct"`
	ConsolidateLevel     uint8  `gorm:"default:0;column:consolidate_level"`
	Exp                  uint32 `gorm:"default:0;column:exp"`
	Level                uint32 `gorm:"default:1;column:level"`
	MAxDmgIn             uint16 `gorm:"default:0;column:max_dmg_in"`
	ArmorIn              uint16 `gorm:"default:0;column:armor_in"`
	EquipAddAtt          []byte `gorm:"column:equip_add_att"`
	TalentPoint          uint8  `gorm:"default:0;column:talent_point"`
	MaxTalentPoint       uint8  `gorm:"default:0;column:max_talent_point"`
	SkillList            []byte `gorm:"column:skill_list"`
	Rating               uint32 `gorm:"default:0;column:rating"`
	ConsolidateLevelStar uint8  `gorm:"default:0;column:consolidate_level_star"`
	AddTalentPoint       int8   `gorm:"default:0;column:add_talent_point"`
	EquipRelAtt          []byte `gorm:"column:equip_rel_att"`
	EquipExAtt           []byte `gorm:"column:equip_ex_att"`
}

func HandleEquip(db1, db2 *gorm.DB) error {

	return nil
}
