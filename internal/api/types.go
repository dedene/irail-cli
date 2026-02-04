package api

// Station represents a Belgian railway station.
type Station struct {
	ID           string `json:"id"`
	URI          string `json:"@id"`
	Name         string `json:"name"`
	StandardName string `json:"standardname"`
	LocationX    string `json:"locationX"`
	LocationY    string `json:"locationY"`
}

// StationsResponse is the response from /stations endpoint.
type StationsResponse struct {
	Version   string    `json:"version"`
	Timestamp string    `json:"timestamp"`
	Station   []Station `json:"station"`
}

// Departure represents a train departure from a station.
type Departure struct {
	ID                  string      `json:"id"`
	Delay               string      `json:"delay"`
	Station             string      `json:"station"`
	StationInfo         Station     `json:"stationinfo"`
	Time                string      `json:"time"`
	Vehicle             string      `json:"vehicle"`
	VehicleInfo         VehicleInfo `json:"vehicleinfo"`
	Platform            string      `json:"platform"`
	PlatformInfo        Platform    `json:"platforminfo"`
	Canceled            string      `json:"canceled"`
	Left                string      `json:"left"`
	IsExtra             string      `json:"isExtra"`
	DepartureConnection string      `json:"departureConnection"`
	Occupancy           Occupancy   `json:"occupancy"`
}

// Arrival represents a train arrival at a station.
type Arrival struct {
	ID                string      `json:"id"`
	Delay             string      `json:"delay"`
	Station           string      `json:"station"`
	StationInfo       Station     `json:"stationinfo"`
	Time              string      `json:"time"`
	Vehicle           string      `json:"vehicle"`
	VehicleInfo       VehicleInfo `json:"vehicleinfo"`
	Platform          string      `json:"platform"`
	PlatformInfo      Platform    `json:"platforminfo"`
	Canceled          string      `json:"canceled"`
	Arrived           string      `json:"arrived"`
	IsExtra           string      `json:"isExtra"`
	ArrivalConnection string      `json:"arrivalConnection"`
}

// Platform contains platform information.
type Platform struct {
	Name   string `json:"name"`
	Normal string `json:"normal"`
}

// VehicleInfo contains vehicle/train information.
type VehicleInfo struct {
	Name      string `json:"name"`
	ShortName string `json:"shortname"`
	Number    string `json:"number"`
	Type      string `json:"type"`
	URI       string `json:"@id"`
	LocationX string `json:"locationX"`
	LocationY string `json:"locationY"`
}

// Occupancy indicates expected train occupancy.
type Occupancy struct {
	URI  string `json:"@id"`
	Name string `json:"name"`
}

// DeparturesContainer wraps the departures array.
type DeparturesContainer struct {
	Number    string      `json:"number"`
	Departure []Departure `json:"departure"`
}

// ArrivalsContainer wraps the arrivals array.
type ArrivalsContainer struct {
	Number  string    `json:"number"`
	Arrival []Arrival `json:"arrival"`
}

// LiveboardResponse is the response from /liveboard endpoint.
type LiveboardResponse struct {
	Version     string              `json:"version"`
	Timestamp   string              `json:"timestamp"`
	Station     string              `json:"station"`
	StationInfo Station             `json:"stationinfo"`
	Departures  DeparturesContainer `json:"departures"`
	Arrivals    ArrivalsContainer   `json:"arrivals"`
}

// Connection represents a journey between two stations.
type Connection struct {
	ID        string          `json:"id"`
	Departure ConnectionStop  `json:"departure"`
	Arrival   ConnectionStop  `json:"arrival"`
	Duration  string          `json:"duration"`
	Vias      ViasContainer   `json:"vias"`
	Occupancy Occupancy       `json:"occupancy"`
	Alerts    AlertsContainer `json:"alerts"`
}

// ConnectionStop represents a departure or arrival stop in a connection.
type ConnectionStop struct {
	Delay               string      `json:"delay"`
	Station             string      `json:"station"`
	StationInfo         Station     `json:"stationinfo"`
	Time                string      `json:"time"`
	Vehicle             string      `json:"vehicle"`
	VehicleInfo         VehicleInfo `json:"vehicleinfo"`
	Platform            string      `json:"platform"`
	PlatformInfo        Platform    `json:"platforminfo"`
	Canceled            string      `json:"canceled"`
	Direction           Direction   `json:"direction"`
	Left                string      `json:"left"`
	Arrived             string      `json:"arrived"`
	Walking             string      `json:"walking"`
	DepartureConnection string      `json:"departureConnection"`
}

// Direction indicates the train's direction.
type Direction struct {
	Name string `json:"name"`
}

// Via represents a transfer point in a connection.
type Via struct {
	ID          string         `json:"id"`
	Arrival     ConnectionStop `json:"arrival"`
	Departure   ConnectionStop `json:"departure"`
	TimeBetween string         `json:"timeBetween"`
	Station     string         `json:"station"`
	StationInfo Station        `json:"stationinfo"`
	Vehicle     string         `json:"vehicle"`
	VehicleInfo VehicleInfo    `json:"vehicleinfo"`
	Direction   Direction      `json:"direction"`
}

// ViasContainer wraps the vias array.
type ViasContainer struct {
	Number string `json:"number"`
	Via    []Via  `json:"via"`
}

// ConnectionsResponse is the response from /connections endpoint.
type ConnectionsResponse struct {
	Version    string       `json:"version"`
	Timestamp  string       `json:"timestamp"`
	Connection []Connection `json:"connection"`
}

// Stop represents a stop on a vehicle's route.
type Stop struct {
	ID                     string   `json:"id"`
	Station                string   `json:"station"`
	StationInfo            Station  `json:"stationinfo"`
	Time                   string   `json:"time"`
	Delay                  string   `json:"delay"`
	Platform               string   `json:"platform"`
	PlatformInfo           Platform `json:"platforminfo"`
	Canceled               string   `json:"canceled"`
	DepartureDelay         string   `json:"departureDelay"`
	ArrivalDelay           string   `json:"arrivalDelay"`
	DepartureCanceled      string   `json:"departureCanceled"`
	ArrivalCanceled        string   `json:"arrivalCanceled"`
	Left                   string   `json:"left"`
	Arrived                string   `json:"arrived"`
	IsExtraStop            string   `json:"isExtraStop"`
	ScheduledDepartureTime string   `json:"scheduledDepartureTime"`
	ScheduledArrivalTime   string   `json:"scheduledArrivalTime"`
}

// StopsContainer wraps the stops array.
type StopsContainer struct {
	Number string `json:"number"`
	Stop   []Stop `json:"stop"`
}

// VehicleResponse is the response from /vehicle endpoint.
type VehicleResponse struct {
	Version     string         `json:"version"`
	Timestamp   string         `json:"timestamp"`
	Vehicle     string         `json:"vehicle"`
	VehicleInfo VehicleInfo    `json:"vehicleinfo"`
	Stops       StopsContainer `json:"stops"`
}

// Composition represents train composition data.
type Composition struct {
	Segments []Segment `json:"segment"`
}

// Segment represents a portion of the train's journey.
type Segment struct {
	Origin      Station            `json:"origin"`
	Destination Station            `json:"destination"`
	Composition SegmentComposition `json:"composition"`
}

// SegmentComposition contains the actual train units.
type SegmentComposition struct {
	Source string `json:"source"`
	Units  Units  `json:"unit"`
}

// Units wraps the unit array.
type Units struct {
	Number string `json:"number"`
	Unit   []Unit `json:"unit"`
}

// Unit represents a single carriage/car.
type Unit struct {
	ID                         string       `json:"id"`
	MaterialType               MaterialType `json:"materialType"`
	HasToilets                 string       `json:"hasToilets"`
	HasAirco                   string       `json:"hasAirco"`
	HasHeating                 string       `json:"hasHeating"`
	HasBikeSection             string       `json:"hasBikeSection"`
	HasPrmSection              string       `json:"hasPrmSection"`
	SeatsFirstClass            string       `json:"seatsFirstClass"`
	SeatsSecondClass           string       `json:"seatsSecondClass"`
	LengthInMeter              string       `json:"lengthInMeter"`
	TractionType               string       `json:"tractionType"`
	CanPassToNextUnit          string       `json:"canPassToNextUnit"`
	TractionPosition           string       `json:"tractionPosition"`
	HasPrmToilet               string       `json:"hasPrmToilet"`
	HasTables                  string       `json:"hasTables"`
	HasSecondClassOutlets      string       `json:"hasSecondClassOutlets"`
	HasFirstClassOutlets       string       `json:"hasFirstClassOutlets"`
	HasSemiPrivateCompartments string       `json:"hasSemiPrivateCompartments"`
}

// MaterialType indicates the type of rolling stock.
type MaterialType struct {
	ParentType  string `json:"parent_type"`
	SubType     string `json:"sub_type"`
	Orientation string `json:"orientation"`
}

// CompositionResponse is the response from /composition endpoint.
type CompositionResponse struct {
	Version     string      `json:"version"`
	Timestamp   string      `json:"timestamp"`
	Composition Composition `json:"composition"`
}

// Alert represents a disturbance or planned work.
type Alert struct {
	ID          string `json:"id"`
	Header      string `json:"header"`
	Description string `json:"description"`
	Lead        string `json:"lead"`
	Link        string `json:"link"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}

// AlertsContainer wraps the alert array.
type AlertsContainer struct {
	Number string  `json:"number"`
	Alert  []Alert `json:"alert"`
}

// Disturbance represents a service disruption.
type Disturbance struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Type        string `json:"type"`
	Timestamp   string `json:"timestamp"`
	Attachment  string `json:"attachment"`
}

// DisturbancesResponse is the response from /disturbances endpoint.
type DisturbancesResponse struct {
	Version     string        `json:"version"`
	Timestamp   string        `json:"timestamp"`
	Disturbance []Disturbance `json:"disturbance"`
}
