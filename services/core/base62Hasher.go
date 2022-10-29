package core

import (
    "strings"
)

var base10To62Map = map[uint64]string {
    0:"0",
    1:"1",
    2:"2",
    3:"3",
    4:"4",
    5:"5",
    6:"6",
    7:"7",
    8:"8",
    9:"9",
    10:"a",
    11:"b",
    12:"c",
    13:"d",
    14:"e",
    15:"f",
    16:"g",
    17:"h",
    18:"i",
    19:"j",
    20:"k",
    21:"l",
    22:"m",
    23:"n",
    24:"o",
    25:"p",
    26:"q",
    27:"r",
    28:"s",
    29:"t",
    30:"u",
    31:"v",
    32:"w",
    33:"x",
    34:"y",
    35:"z",
    36:"A",
    37:"B",
    38:"C",
    39:"D",
    40:"E",
    41:"F",
    42:"G",
    43:"H",
    44:"I",
    45:"J",
    46:"K",
    47:"L",
    48:"M",
    49:"N",
    50:"O",
    51:"P",
    52:"Q",
    53:"R",
    54:"S",
    55:"T",
    56:"U",
    57:"V",
    58:"W",
    59:"X",
    60:"Y",
    61:"Z",
}

func Base62hash(base10Number uint64) (string, error) {

    quotient := base10Number
    reminder := base10Number
    reminders := []uint64{}

    for ok := true; ok; ok = quotient != 0 {

        reminder = quotient % 62
        quotient = quotient / 62
        reminders = append(reminders, reminder)
    }

    var hashValue strings.Builder
    for i := len(reminders)-1; i >= 0; i-- {
        hashValue.WriteString(base10To62Map[reminders[i]])
    }

    return hashValue.String(), nil
}
