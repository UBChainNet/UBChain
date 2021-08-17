package exchange_runner

var exMethods = map[string]*MethodInfo{
	"Methods": &MethodInfo{
		Name:   "Methods",
		Params: nil,
		Returns: []Value{
			{
				Name: "Open methods",
				Type: "",
			},
		},
	},
	"MethodExist": &MethodInfo{
		Name: "MethodExist",
		Params: []Value{
			{
				Name: "method",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "exist",
				Type: "bool",
			},
		},
	},
	"Pairs": &MethodInfo{
		Name:   "Pairs",
		Params: nil,
		Returns: []Value{
			{
				Name: "pair list",
				Type: "",
			},
		},
	},
	"ExchangeRouter": &MethodInfo{
		Name: "ExchangeRouter",
		Params: []Value{
			{
				Name: "tokenA",
				Type: "string",
			},
			{
				Name: "tokenB",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "paths",
				Type: "",
			},
		},
	},
}

var pairMethods = map[string]*MethodInfo{
	"Methods": &MethodInfo{
		Name:   "Methods",
		Params: nil,
		Returns: []Value{
			{
				Name: "Open methods",
				Type: "",
			},
		},
	},
	"MethodExist": &MethodInfo{
		Name: "MethodExist",
		Params: []Value{
			{
				Name: "method",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "exist",
				Type: "bool",
			},
		},
	},
}
