// Pattern 1

// model

type User struct {
  Name string  `json:"name"`  validate:"nonzero"
  Age  uint    `json:"age"`   validate:"min=1"
  Address string `json:"address"` validate:"nonzero"
}

// unmarshalling json byte slice to object
var user User
if err := json.NewDecoder(jsonByteSlice).Decode(&user); err != nil {...}

// marshalling  object to json byte slice
if jsonByteSlice, err := json.Marshal(object); err != nil {...}

// validation

if err := validator.Validate(user); err != nil {...}

// Pattern 2

type User struct {
     Name          *string `json:"name"`              // required, but no defaults
     Age           *uint   `json:"age,omitempty"`     // optional
     Address       *string `json:"address,omitempty"` // optional
     FavoriteColor string  `json:"favoriteColor"`     // required, uses defaults
}

// unmarshalling
var user User
if err := json.NewDecoder(jsonByteSlice).Decode(&user); err != nil {...}

// marshalling
if jsonByteSlice, err := json.Marshal(object); err != nil {...}

// validation
func Validate(user User) {
     // default - validate value
     // optional - if non nil, validate value
     // required non default - validate not nil, then validate value
}

// Pattern 3

UserId *uint `validate:"nonzero,min=100"`

type User struct {
     Name          *string `json:"name"              validate:"nonzero,min=1"` // required, but no defaults
     Age           *uint   `json:"age,omitempty"     validate:"min=1"`         // optional
     Address       *string `json:"address,omitempty" validate:"min=1"`         // optional
     FavoriteColor string  `json:"favoriteColor"`                              // required, uses defaults
}

// unmarshalling
var user User
if err := json.NewDecoder(jsonByteSlice).Decode(&user); err != nil {..}

// marshalling
if jsonByteSlice, err := json.Marshal(object); err != nil {...}

// validation
if errs := validator.Validate(user); errs != nil {...}


// Pattern 4
type UserPostRequest struct {
     Name               *string `json:"name" validate:"nonzero"` // required, but no defaults
     Age                *uint   `json:"age"`                     // optional
     Address            *string `json:"address"`                 // optional
     FavoriteInstrument string  `json:"favoriteInstrument"`      // required, uses defaults
     FavoriteColor      *string `json:"favoriteColor"`           // required, uses custom defaults
}

// model
type User struct {
     Name               string  `json:"name" validate:"nonzero"`            // required
     Age                *uint   `json:"age,omitempty" validate:"min=1"`     // optional
     Address            *string `json:"address,omitempty" validate:"min=1"` // optional
     FavoriteInstrument string  `json:"favoriteInstrument"`                 // required
     FavoriteColor      string  `json:"favoriteColor" validate:"nonzero"`   // required
}

// unmarshalling
var postRequest UserPostRequest
if err := json.NewDecoder(jsonByteSlice).Decode(&postRequest); err != nil {..}
if errs := validator.Validate(postRequest); errs != nil {...}

user.Name = *postRequest.Name
user.Age = postRequest.Age
user.Address = postRequest.Address
user.FavoriteInstrument = postRequest.FavoriteInstrument
user.FavoriteColor = "blue"
if postRequest.FavoriteColor != nil {
  user.FavoriteColor = *postRequest.FavoriteColor
}
if errs := validator.Validate(user); errs != nil {...}

// marshalling
if jsonByteSlice, err := marshal(object); err != nil {...}

// validation
if errs := validator.Validate(user); errs != nil {...}