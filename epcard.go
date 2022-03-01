package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "image"
    "image/color"
    "image/png"
    "log"
    "net/http"

    "github.com/fogleman/gg"
)

type Data struct {
    BgImgPath string
    FontPath  string
    is_trash  bool
    owned     bool
}

type Texts struct {
    Catching string
    Price    string
    Tax      string
    bonus    string
    money    string
    name     string
    detail   string
    place    string
}

func GenEpcard(data Data, texts Texts) (image.Image, error) {
    bgImage, err := gg.LoadImage(data.BgImgPath)
    if err != nil {
        return nil, err
    }
    imgWidth := bgImage.Bounds().Dx()
    imgHeight := bgImage.Bounds().Dy()

    dc := gg.NewContext(imgWidth, imgHeight)
    dc.DrawImage(bgImage, 0, 0)

    maxWidth := float64(imgWidth)
    if err := dc.LoadFontFace(data.FontPath, 24); err != nil {
        return nil, err
    }

    //===================================
    dc.LoadFontFace(data.FontPath, 24)
    dc.SetColor(color.White)
    dc.DrawStringWrapped(texts.Catching, 311, 35, 0.5, 1, maxWidth, 1, gg.AlignLeft)
    if !data.is_trash {
        dc.SetColor(color.Black)
        dc.DrawStringWrapped(texts.money, 325, 173, 0.5, 1, maxWidth, 1, gg.AlignLeft)
    }

    dc.LoadFontFace(data.FontPath, 20)
    dc.SetColor(color.Black)
    dc.DrawStringWrapped(texts.name, 590, 89, 0.5, 1, maxWidth, 1, gg.AlignLeft)
    dc.DrawStringWrapped(texts.Price, 320, 83, 0.5, 1, maxWidth, 1, gg.AlignLeft)

    if !data.is_trash {
        dc.LoadFontFace(data.FontPath, 16)
        dc.SetHexColor("320000")
        dc.DrawStringWrapped(texts.Tax, 312, 104, 0.5, 1, maxWidth, 1, gg.AlignLeft)
        dc.SetHexColor("003200")
        dc.DrawStringWrapped(texts.bonus, 312, 122, 0.5, 1, maxWidth, 1, gg.AlignLeft)
    }

    dc.LoadFontFace(data.FontPath, 14)
    dc.SetColor(color.Black)
    dc.DrawStringWrapped(texts.detail, 590, 154, 0.5, 1, maxWidth, 1.5, gg.AlignLeft)

    dc.LoadFontFace(data.FontPath, 8)
    dc.DrawStringWrapped(texts.place, 530, 178, 0.5, 1, maxWidth, 1.5, gg.AlignLeft)

    if data.owned {
        dc.LoadFontFace(data.FontPath, 10)
        dc.SetColor(color.White)
        dc.DrawStringWrapped("낚시터 주인", 312, 49, 0.5, 1, maxWidth, 1, gg.AlignLeft)
    }

    return dc.Image(), nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Rarity string `json:"rarity"`
        Trash  bool   `json:"is_trash"`
        Owned  bool   `json:"owned"`
        Catch  string `json:"catching"`
        Price  string `json:"price"`
        Tax    string `json:"tax"`
        Bonus  string `json:"bonus"`
        Money  string `json:"money"`
        Name   string `json:"name"`
        Detail string `json:"detail"`
        Place  string `json:"place"`
    }
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    Bgpath := fmt.Sprintf("./assets/themes/default/rank-%s.png", req.Rarity)
    if req.Rarity == "" {
        Bgpath = "./assets/themes/default/default.png"
    }

    img, err := GenEpcard(
        Data{
            BgImgPath: Bgpath,
            FontPath:  "./assets/SpoqaHanSansNeo-Bold.ttf",
            is_trash:  req.Trash,
            owned:     req.Owned,
        },
        Texts{
            Catching: req.Catch,
            Price:    req.Price,
            Tax:      req.Tax,
            bonus:    req.Bonus,
            money:    req.Money,
            name:     req.Name,
            detail:   req.Detail,
            place:    req.Place,
        },
    )
    if err != nil {
        log.Println(err)
    }

    buffer := new(bytes.Buffer)
    if err := png.Encode(buffer, img); err != nil {
        log.Println("unable to encode image")
    }

    w.Header().Set("Content-Type", "image/png")
    w.Write(buffer.Bytes())
}

func main() {
    handler := http.HandlerFunc(handleRequest)
    http.Handle("/card", handler)
    fmt.Println("Server started at port 5050")
    http.ListenAndServe(":5050", nil)
}