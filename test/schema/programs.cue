package schemas

// # The Anatomy of a Program.
//
#program: {
	// ## Name
	//
	// A human readable name/label for the program.
	//
	//   - Usually this will be what the client calls the program in a "nicer"
	//     format (e.g. "2022 HERBICIDE GROWTH" might become
	//     "2022 Retail → Herbicide Growth")
	//   - Use ` → ` as a logical separator
	//   - Use ® and ™ where appropriate when including a progduct/brand name
	//   - Consider other programs in the program-set, try to ensure similar
	//     programs/offers/components sort together
	name: string

	// ## Program
	//
	// A machine readable key for the program.
	//
	//   - Should be composed of components of the label
	//     (e.g. `2022.herbicide_growth`)
	//   - Only use lower-case alpha characters, numbers, `_` or `.`
	//   - A period (`.`) should be used when separating "components" of a key
	//     (e.g. `2022.herbicide...`)
	//   - An underscore (`_`) should be sued when separating words in a component
	//     (e.g. `herbicide_growth`)
	//   - When multiple programs configuration exist for a similar purpose,
	//     compose a key that begins with the "widest" component and ends with the
	//     "narrowest" (e.g. `2022.herbicide_growth.southeast.level_1`)
	//   - Consider other programs in the program-set, try to ensure similar
	//     programs/offers/components sort together
	program: string

	// ## Aggregate Program Key
	//
	// A machine readable key for grouping program configurations together.
	//
	//   - Optional, defaults to `program` when not provided
	//   - Use when multiple configureations are required to represent on logical
	//     program (e.g. `program: 2022.base.level_1` might have
	//     `aggregate_program_key: `2022.base`)
	aggregate_program_key: *program | string

	// Strategy
	//
	// A key representing the approach Incent will take to process
	//   transactions
	//
	//   - Not explicitly set.
	//   - See `./strategies.cue` for further detail
	strategy: string

	// Calculator
	//
	// A key representing how Incent will calculate potential earnings
	//   - note: not all calculators are compatible with every strategy
	calculator: "allow_any_product" | "annualized_fixed_amount_per_purchased_unit" | "annualized_percent_of_purchased_amount" | "fixed_amount" | "fixed_amount_per_purchased_unit" | "fixed_amount_per_treated_acre" | "percent_of_incentive_amount" | "percent_of_purchased_amount" | "percent_of_amount" | "fixed_amount_per_configured_value"

	// ## Safe-Mode
	//
	// Indicates whether the generated potential earnings should be processed in
	//   "safe mode" or not.
	//
	//   - Will return only positive incentive values regardless of setting.
	//   - When `true` (default), net-negative volumes/amounts will _not_ for one
	//     sku will not impact other skus.
	//   - When `false`, a net negative of one sku is allowed to impact positive
	//     purchases of another product within a program.
	//
	//   This does not solve when there are multiple offers in a program-set that
	//   impact each other as the resulting response for a net-negative offer will
	//   be $0. For example, if there are three offers that aggregate up to one
	//   offer and one of the offers results in a net negative earning then the
	//   overall incentive for the payment group needs to be reduced by the net
	//   negative amount. Safe-mode does not look across programs.
	safe_mode?: *true | bool

	// ## Allow Negative Earnings
	//
	// Indicates where the generated potential earnings should be processed in
	// "allow negative earnings" mode or not.
	//
	//   - When `false` (default) see `safe_mode` for how potential earning will
	//     be generated
	//   - When `true` the program is automatically set to "unsafe-mode"
	//     (e.g. `safe_mode: false`)
	//   - When `true` net-negative values for one program will be allowed to
	//     impact another program.
	//
	//  `allow_negative_earnings` should only be used to handle situations where
	//  we qualify at one level and earn at another level so that proration can be
	//  handled externally. We plan to implement a new feature that will
	//  ultimately replace allow negative earnings so we do not want to implement
	//  new programs with this feature unless absolutely necessary. Allow for
	//  negative calculation results when calculating for multiple locations in
	//  one request (purposes of proration).
	allow_negative_earnings?: *false | bool

	// ## Dynamic
	//
	// This program is intended to be used as a "Custom Offer"
	//
	//   - Cutomer offers contain very little or no business rules
	//   - Include required attributes and metadata
	dynamic?: *false | bool

	// ## Variable Rate
	//
	// Whether or not the program expects incentive rates provided in the request.
	//
	//   - DEPRECATED: new offers should _not_ be configured with this setting
	//   - Supersceded by `dynamic`
	variable_rate?: *false | bool

	// ## Applies To
	//
	// Indicates what part of the calculation pipeline the program is intended to
	// be used for.
	//
	//   - `pre_processing`: DEPRECATED: see `exclude_from_processing`. was for
	//     use when the the program will be used as part of.
	//     the append-properties pre-processor.
	//   - `transactions`: (default) traditional usage; calculate potential
	//     earnings based on the provided transactions.
	//   - `incentives`: used when a program applies to potential earnings
	//     calculated in a traditional program. e.g. earn 2% over base earnings
	applies_to?: ["pre_processing" | "transactions" | "incentives"]
	// ## Exclude from Processing
	//
	//  Used when a program is configured expressly for use in an
	//  append-properties pre-processor.
	//
	//   - Superscedes `applies_to.pre_processing`
	exclude_from_processing?: *false | bool

	// ## Description
	//
	// A sentence or two describing the program.
	//
	//   - optional: defaults to the program name/label
	//   - Typically verbiage from actual program sheet authored by client.
	description: *name | string

	// Metadata
	//
	// Key value pairs containing informational data not directly impacting
	// program calculation.
	//
	//   - This data is supplied to users (requestors) in various ways via the
	//     API.
	//   - Often includes program labls/keys/codes that relate back to the actual
	//     program as described by the client.
	//   - Can include a list of required payment-group properties.
	//   - Can include a list of required transaction configs describing the
	//     dataset the program expects (for qualification and incentive
	//     calculation).
	metadata?: [...]

	// ## Lookups
	//
	// Program details that are used for validation and elsewhere in
	//   qualifying and/or calculation.
	//
	//   - geography: see ./common.cue:#geography
	//   - chem_products: see ./common.cue:#product_struct
	lookups: {
		geography:     #geography
		chem_products: #products_struct
	}

	// ## Qualifications
	//
	// Describe business rules included in the program that are required to be met
	// in order for a payment-group to "qualify" for earning.
	//
	//  - see [qualifications.cue](qualifications.html#attribute-#maximum_total_chem_amount_qualification)
	//    for further details
	qualifications: {...}

	// ## Eligibility
	//
	// An expression based on qualification keys that is used to "select" programs
	// a payment-group is eligible to earn in.
	//
	//   - Most commonly used for `geography`, `context`, or other "property"
	//     based qualfications.
	//   - Assuming two configured qualifications named `geographical_eligibility`
	//     and `context_eligibility` that are _both required_ the eligibility
	//     would be simply:
	//     ```
	//     geographical_eligibility AND context_eligibility
	//     ```
	//   - Eligibility is intended as an extremely fast way to chose programs for
	//     a given request.
	//   - The eligibility expression should _never_ contain the key of a
	//     qualification that uses transaction filtering (e.g. `filters`)
	//   - Incent uses the library Dentaku to evaluate the given expression. See
	//     here for further details:
	//     https://github.com/rubysolo/dentaku#built-in-operators-and-functions
	eligibility?: string | bool

	// ## Qualifier
	//
	// An expression based on qualification keys that is used to determine whether
	// or not the payment-group is "qualified" to earn in the program.
	//
	//   - Assuming two configured qualifications named `2022_edi_qualification`
	//     and `2022_2021_edi_yoy_qualification` where meeting _either_ threshold
	//     would allow the payment-group to qualify, the qualifier would be:
	//     ```
	//     2022_edi_qualification OR 2022_2021_edi_yoy_qualification
	//     ```
	//   - More complex expression can be constructed using `NOT` and/or
	//     parentheses. For example (with abbreviated keys):
	//     ```
	//     (a AND b) OR (NOT(a) AND c)
	//     ```
	//   - Incent uses the library Dentaku to evaluate the given expression. See
	//     here for further details:
	//     https://github.com/rubysolo/dentaku#built-in-operators-and-functions
	//   - See [qualifications](#attribute-qualifications)
	qualifier: string | bool

	...
}
