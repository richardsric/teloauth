package public

type message struct {
	Sucess      string
	Error       string
	Warning     string
	Information string
}
type pageData struct {
	Code      string
	TpOWnerID string
}
type tradeProfileData struct {
	TradeProfile           string
	TradeType              int64
	TradeMode              string
	StopLoss               float64
	ProfitLockStart        float64
	WalletExposure         float64
	ExchangeID             int64
	BuyOrderTimeout        int64 // not fixed
	ProfitKeep             float64
	SellTrigger            float64
	InheritSubscribersFrom int64
	MinCap                 float64
	MaxCap                 float64
	Market                 string
	BaseMarket             string
	PartialBuyTimeout      int64
	PartialBuyTimeoutPl    float64 // not fixed
	SellOrderTimeout       int64
	ProfileID              int64
	ProfilePrivacy         string
	ProfitKeepReadjustPl   float64
	ProfitKeepReadjust     float64
	SellTriggerReadjust    float64
	TradeCommission        float64
	TpOwnerID              int64
	Code                   string
	InheritSubDesc         string
}

type pageMainData struct {
	Title          string
	TradeProfile   tradeProfileData
	BaseMarketData map[string]string
	TradeTypesData map[int]string
	ExchangeData   map[int]string
	SubscriberTp   map[int]string
}

type signPageData struct {
	Title        string
	AccountTypes map[int]string
	CountryCode  map[string]string
	AccountInfo  accountInfo
}

type accountInfo struct {
	FirstName        string
	LastName         string
	Email            string
	Phone            string
	AccountType      int64
	AllowAPIWithdraw int64
	FacebookID       string
	Whatsapp         string
	CountryCode      string
	State            string
	RegType          int64
	CountryName      string
}

//user is a struct for the retrieved and an authentiacted user.
type user struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

type exchange struct {
	ExchnageID   int
	ExchangeName string
}
type allExchanges struct {
	Exchanges []exchange
}

type tradeType struct {
	TradeTypeID   int
	TradeTypeName string
}
type allTradeTypes struct {
	TradeTypes []tradeType
}

type baseMarket struct {
	baseMarketID   int
	baseMarketName string
}

type allBaseMarket struct {
	BaseMarkets []baseMarket
}
type allDropDownData struct {
	Exchanges   allExchanges
	TradeTypes  allTradeTypes
	BaseMarkets allBaseMarket
}
