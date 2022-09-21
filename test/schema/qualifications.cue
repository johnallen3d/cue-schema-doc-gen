package schemas

#value_threshold: {
	type:  "value"
	value: number | string
}

#property_threshold: {
	type:  "key"
	value: string | bool
}

#expression_threshold: {
	type:  "expression"
	value: string
}

#threshold: #value_threshold | #property_threshold | #expression_threshold

#geographical_eligibility: {
	type: "geography" | *"geography"
}

#qualification_optional_fields: {
	label?:     string
	waive_key?: string
	variance?: {
		measure: "percent" | number
		value:   number | string
	}
	exclusive: bool | *false
}

//minimum qualifications

#min_qualification: {
	minimum: #threshold
	filters: #filters
}

#minimum_chem_volume_qualification: {
	#min_qualification
	#qualification_optional_fields

	type: "minimum_chem_volume" | *"minimum_chem_volume"
}

#minimum_configured_value_qualification: {
	#qualification_optional_fields

	type:    "minimum_configured_value" | *"minimum_configured_value"
	minimum: #threshold
	key:     string
}

#minimum_total_chem_amount_qualification: {
	#min_qualification
	#qualification_optional_fields

	type: "minimum_total_chem_amount" | *"minimum_total_chem_amount"
}

#minimum_total_chem_volume_qualification: {
	#min_qualification
	#qualification_optional_fields

	type: "minimum_total_chem_volume" | *"minimum_total_chem_volume"
}

#minimum_total_configured_count_qualification: {
	#min_qualification
	#qualification_optional_fields

	type: "minimum_total_configured_count" | *"minimum_total_configured_count"
	key:  string
}

#minimum_total_incentive_result_qualification: {
	#qualification_optional_fields

	type:     "minimum_total_incentive_result" | *"minimum_total_incentive_result"
	minimum:  #threshold
	filters?: #filters
}

#minimum_total_product_count_qualification: {
	#min_qualification
	#qualification_optional_fields

	type: "minimum_total_product_count" | *"minimum_total_product_count"
}

#minimum_total_product_coverage_qualification: {
	#min_qualification
	#qualification_optional_fields

	type: "minimum_total_product_coverage" | *"minimum_total_product_coverage"
}

//maximum qualifications

#maximum_total_chem_amount_qualification: {
	#qualification_optional_fields

	type:    "maximum_total_chem_amount" | *"maximum_total_chem_amount"
	maximum: #threshold
	filters: #filters
}

//property qualifications

#property_qualification: {
	type:       "property" | *"property"
	key:        string
	value:      string | bool
	waive_key?: string
	label?:     string
}

#property_date_range_qualification: {
	type:       "property_date_range" | *"property_date_range"
	key:        string
	date_range: #date_range_filter
}

//ratio qualifications

#minimum_configured_ratio_qualification: {
	#qualification_optional_fields

	type:                 "minimum_configured_ratio" | *"minimum_configured_ratio"
	minimum:              #threshold
	accept_zero_by_zero?: bool
	numerator:            #ratio_fraction_part | #filter_fraction_part | #property_fraction_part
	denominator:          #ratio_fraction_part | #filter_fraction_part | #property_fraction_part
}

#minimum_total_chem_acres_ratio_qualification: {
	type:        "minimum_total_chem_acres_ratio" | *"minimum_total_chem_acres_ratio"
	minimum:     #threshold
	numerator:   #filter_fraction_part
	denominator: #filter_fraction_part
	waive_key?:  string
	label?:      string
}

#range_configured_ratio_qualification: {
	#qualification_optional_fields

	type: "range_configured_ratio" | *"range_configured_ratio"
	range: {
		minimum: #threshold
		maximum: #threshold
	}
	accept_zero_by_zero?: bool
	numerator:            #ratio_fraction_part | #filter_fraction_part | #property_fraction_part
	denominator:          #ratio_fraction_part | #filter_fraction_part | #property_fraction_part
}
