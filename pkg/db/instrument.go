//=============================================================================
/*
Copyright © 2023 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package db

import "gorm.io/gorm"

//=============================================================================

func GetInstruments(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]Instrument, error) {
	var list []Instrument
	res := tx.Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, res.Error
	}

	return &list, nil
}

//=============================================================================

func GetInstrumentsAsMap(tx *gorm.DB) (map[uint]*Instrument, error) {
	var list []Instrument
	res := tx.Find(&list)

	if res.Error != nil {
		return nil, res.Error
	}

	instMap := map[uint]*Instrument{}

	for _, inst := range list {
		inAux := inst
		instMap[inst.Id] = &inAux
	}

	return instMap, nil
}

//=============================================================================

func GetOrCreateInstrument(tx *gorm.DB, ticker string, i *Instrument) (*Instrument, error) {
	res := tx.Where(&Instrument{Ticker: ticker}).FirstOrCreate(i)

	if res.Error != nil {
		return nil, res.Error
	}

	return i, nil
}

//=============================================================================
