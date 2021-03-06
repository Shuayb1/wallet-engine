package models

import (
	"errors"
	"html"
	"strings"
	"time"
	"fmt"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Wallet struct {
	WalletID        uint32    `gorm:"primary_key;auto_increment" json:"wallet_id"`
	WalletBalance   uint32    `json:"walletBalance"`
	Nickname  string    `gorm:"size:255;not null;unique" json:"nickname"`
	Status  string    `gorm:"size:255;not null;unique" json:"status"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type CreditDebit struct {
	CreditAmount uint32 `json:"credit_amount"`
	DebitAmount uint32 `json:"debit_amount"`

}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *Wallet) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *Wallet) Prepare() {
	u.WalletID = 0
	u.WalletBalance = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Status = "Active"
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}



func (u *Wallet) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("Required Wallet ID")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *Wallet) SaveWallet(db *gorm.DB) (*Wallet, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &Wallet{}, err
	}
	return u, nil
}

func (u *Wallet) FindAllWallets(db *gorm.DB) (*[]Wallet, error) {
	var err error
	wallets := []Wallet{}
	err = db.Debug().Model(&Wallet{}).Limit(100).Find(&wallets).Error
	if err != nil {
		return &[]Wallet{}, err
	}
	return &wallets, err
}

func (u *Wallet) FindWalletByID(db *gorm.DB, uid uint32) (*Wallet, error) {
	var err error
	err = db.Debug().Model(Wallet{}).Where("wallet_id = ?", uid).Take(&u).Error
	if err != nil {
		return &Wallet{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Wallet{}, errors.New("Wallet Not Found")
	}
	return u, err
}

func (u *Wallet) CreditAWallet(db *gorm.DB, uid uint32, amount uint32) (*Wallet, error) {

	
	type Result struct {
		Walletbalance  uint32
		Status string
	}

	var result Result


	getbalanceAndStatus := db.Debug().Model(&Wallet{}).Select("wallet_balance, status").Where("wallet_id = ?", uid).Scan(&result)
	fmt.Println("CreditAWallet: model getbalanceAndStatus",getbalanceAndStatus)

	if getbalanceAndStatus.Error != nil {
		return &Wallet{}, db.Error
	}

	fmt.Println("CreditAWallet: model result",result)


	newbalance := result.Walletbalance + amount
	status := result.Status
	fmt.Println("CreditAWallet: model status",amount,result.Walletbalance, newbalance, status)


	if status != "Active"{

		return u, errors.New("Account is inactive")
	}

	db = db.Debug().Model(&Wallet{}).Where("wallet_id=?", uid).Update("wallet_balance", newbalance)
	if db.Error != nil {
		return &Wallet{}, db.Error
	}

	// This is the display the updated wallet
	err := db.Debug().Model(&Wallet{}).Where("wallet_id = ?", uid).Take(&u).Error
	fmt.Println("CreditAWallet: model err", err)

	if err != nil {
		return &Wallet{}, err
	}

	fmt.Println("CreditAWallet: model err", err)

	return u, nil
}


func (u *Wallet) DebitAWallet(db *gorm.DB, uid uint32, amount uint32) (*Wallet, error) {

	
	type Result struct {
		Status string
		Walletbalance  uint32
	  }
	  
	var result Result
	getbalanceAndStatus := db.Model(&Wallet{}).Select("wallet_balance, status").Where("wallet_id = ?", uid).Scan(&result)
	
	if getbalanceAndStatus.Error != nil {
		return &Wallet{}, db.Error
	}

	if amount > result.Walletbalance {
		return &Wallet{},errors.New("Insufficient Balance")
	}
	
	newbalance := result.Walletbalance - amount
	status := result.Status

	if status != "Active" {
		return &Wallet{}, errors.New("Account is inactive")
	}

	db = db.Debug().Model(&Wallet{}).Where("wallet_id=?", uid).Update("wallet_balance", newbalance)
	if db.Error != nil {
		return &Wallet{}, db.Error
	}

	// This is the display the updated wallet
	err := db.Debug().Model(&Wallet{}).Where("wallet_id = ?", uid).Take(&u).Error
	if err != nil {
		return &Wallet{}, err
	}
	return u, nil
}

func (u *Wallet) ActivateAWallet(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Wallet{}).Where("wallet_id=?", uid).Update("status", "Active")
	if db.Error != nil {
		return 0, db.Error
	}

	// This is the display the updated wallet
	err := db.Debug().Model(&Wallet{}).Where("wallet_id = ?", uid).Take(&u).Error
	if err != nil {
		return 0, err
	}
	return 0, nil

}

func (u *Wallet) DeactivateAWallet(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Wallet{}).Where("wallet_id=?", uid).Update("status", "Inactive")
	if db.Error != nil {
		return 0, db.Error
	}

	// This is the display the updated wallet
	err := db.Debug().Model(&Wallet{}).Where("wallet_id = ?", uid).Take(&u).Error
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (u *Wallet) DeleteAWallet(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Wallet{}).Where("wallet_id = ?", uid).Take(&Wallet{}).Delete(&Wallet{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
