package restapi

import "slices"

// Define the Rest API interface. You should not modify the code in this block.

type User struct {
	Name    string             `json:name`
	Owes    map[string]float64 `json:owes`
	OwedBy  map[string]float64 `json:owedBy`
	Balance float64            `json:balance`
}

type GetUsersRequest struct {
	Users []string
}

type GetUsersResponse struct {
	Users []User
}

type AddUserRequest struct {
	User string
}

type AddUserResponse struct {
	User User
}

type AddIouRequest struct {
	Lender   string
	Borrower string
	Amount   float64
}

type AddIouResponse struct {
	Users []User
}

type RestApi interface {
	GetUsers(GetUsersRequest) GetUsersResponse
	AddUser(AddUserRequest) AddUserResponse
	AddIou(AddIouRequest) AddIouResponse
}

// Your code goes below here. Implement the RestApi interface.

type Api struct {
	database []User
}

func NewApi(database []User) RestApi {
	return &Api{database: database}
}

func (a *Api) GetUsers(req GetUsersRequest) GetUsersResponse {
	if len(req.Users) == 0 {
		return GetUsersResponse{Users: a.database}
	}

	var filtered []User

	for _, name := range req.Users {
		for _, u := range a.database {
			if u.Name == name {
				filtered = append(filtered, u)
			}
		}
	}
	return GetUsersResponse{Users: filtered}
}

func (a *Api) AddUser(req AddUserRequest) AddUserResponse {
	newUser := User{
		Name:   req.User,
		Owes:   make(map[string]float64),
		OwedBy: make(map[string]float64),
	}

	a.database = append(a.database, newUser)
	return AddUserResponse{User: newUser}
}

func (a *Api) AddIou(req AddIouRequest) AddIouResponse {
    var lender, borrower *User

    for i := range a.database {
        if a.database[i].Name == req.Lender {
            lender = &a.database[i]
        }
        if a.database[i].Name == req.Borrower {
            borrower = &a.database[i]
        }
    }

    if lender != nil && borrower != nil {
        lender.Balance += req.Amount
        borrower.Balance -= req.Amount

        alreadyOwes := lender.Owes[borrower.Name]

        if alreadyOwes > 0 {
            if req.Amount < alreadyOwes {
                lender.Owes[borrower.Name] -= req.Amount
                borrower.OwedBy[lender.Name] -= req.Amount
            } else if req.Amount == alreadyOwes {
                delete(lender.Owes, borrower.Name)
                delete(borrower.OwedBy, lender.Name)
            } else {
                remaining := req.Amount - alreadyOwes
                delete(lender.Owes, borrower.Name)
                delete(borrower.OwedBy, lender.Name)
                
                lender.OwedBy[borrower.Name] += remaining
                borrower.Owes[lender.Name] += remaining
            }
        } else {
            lender.OwedBy[borrower.Name] += req.Amount
            borrower.Owes[lender.Name] += req.Amount
        }
    }

    res := []User{*lender, *borrower}
    slices.SortFunc(res, func(a, b User) int {
        if a.Name < b.Name { return -1 }
        return 1
    })
    return AddIouResponse{Users: res}
}
