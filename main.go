package main

import (
    "encoding/json"
    "net/http"
    "sort"
    "log"
    "strconv"
    "fmt"
    "os"

    "github.com/gorilla/mux"
)


type ContainerMap map[int]Container

type Container struct {
    Volume                float32     `json:"volume"`
    HazardPrLitre         float32     `json:"hazardPrLitre"`
    RemovalCostPrLitre    float32     `json:"removalCostPrLitre"`
}

var Containers ContainerMap


type RemovalStruct struct {
    RemovalCost           float32     `json:"removalCost"`
}

var Removal RemovalStruct


type Cont struct {
    id            int
    vol           float32
    hazardPrLt    float32
    remCostPrLt   float32
}


type LiquidatedMaterial struct {
    Id                  int            `json:"id"`
    VolumeToLiquidate   float32        `json:"volumeToLiquidate"`
}




const PORT int = 8080

const InsuranceConstant = 3


func main() {

    r := mux.NewRouter()

    Containers = make(ContainerMap)

    r.HandleFunc("/container/{id}", AddOrUpdateContainer).Methods("PUT")
    r.HandleFunc("/liquidate", Liquidate).Methods("POST")
    r.HandleFunc("/showcontainers", ShowContainers).Methods("GET")
    r.HandleFunc("/showspare", ShowAvailableRemovalSpare).Methods("GET")

    fmt.Printf("SERVER LISTENING ON :%d\n", PORT)
    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(PORT), r))

}



func ShowContainers(w http.ResponseWriter, r *http.Request) {

     json.NewEncoder(w).Encode(Containers)

}



func ShowAvailableRemovalSpare(w http.ResponseWriter, r *http.Request) {

     json.NewEncoder(w).Encode(Removal)

}



func Liquidate(w http.ResponseWriter, r *http.Request) {

    var RemovalSpare float32
    var rem RemovalStruct
    var consumedSpare float32
    var liquidatedVol float32
    var LiquidatedList = []LiquidatedMaterial{}


    result := json.NewDecoder(r.Body).Decode(&rem)

    if result != nil {
        fmt.Fprintf(os.Stderr, "result=%v\n", result)
        return
    }

    Removal.RemovalCost += rem.RemovalCost

    fmt.Printf("RemovalStruct %v\n", Removal.RemovalCost)

    RemovalSpare = Removal.RemovalCost

    var ContList []Cont

    for contId, container := range Containers {
        fmt.Println("Id:", contId, "Container:", container)
        ContList = append(ContList, Cont{
                      contId,
                      container.Volume,
                      container.HazardPrLitre,
                      container.RemovalCostPrLitre,
                   })
    }

    fmt.Printf("ContList %v\n", ContList)

    sort.SliceStable(ContList, func(i, j int) bool {
        return ContList[i].hazardPrLt > ContList[j].hazardPrLt })

    fmt.Printf("Re-ordered ContList %v\n", ContList)

    for _, c := range ContList {

        fmt.Printf("RemovalSpare %v   RemovalCostPrLitre %v   HazardPrLitre %v   Volume %v    Id %d\n",
                   RemovalSpare, c.remCostPrLt, c.hazardPrLt, c.vol, c.id)

        if int(RemovalSpare / c.remCostPrLt) > 0 {
            if RemovalSpare / c.remCostPrLt >= c.vol {

                liquidatedVol = c.vol
                consumedSpare = c.remCostPrLt * liquidatedVol

                fmt.Printf("Available RemovalSpare %g => now Liquidating [RemovalCostPrLitre x Volume] %g x %g = %g\n",
                           RemovalSpare, c.remCostPrLt, liquidatedVol, consumedSpare)

                RemovalSpare = RemovalSpare - consumedSpare

                delete(Containers, c.id)

                fmt.Printf("Liquidated %g Litres, Consumed Spare costs %g, Remaining Spare costs %g\n",
                           liquidatedVol, consumedSpare, RemovalSpare)

                LiquidatedList = append(LiquidatedList,
                                        LiquidatedMaterial{
                                                   c.id,
                                                   liquidatedVol,
                                        })

            } else {

                liquidatedVol = RemovalSpare / c.remCostPrLt
                consumedSpare = c.remCostPrLt * liquidatedVol

                fmt.Printf("Available RemovalSpare %g => now Liquidating [RemovalCostPrLitre x Volume] %g x %g = %g\n",
                           RemovalSpare, c.remCostPrLt, liquidatedVol, consumedSpare)

                RemovalSpare = RemovalSpare - consumedSpare

                Containers[c.id] = Container{
                                         c.vol - liquidatedVol,
                                         c.hazardPrLt,
                                         c.remCostPrLt,
                                   }

                fmt.Printf("Liquidated %g Litres, Consumed Spare costs %g, Remaining Spare costs %g\n",
                           liquidatedVol, consumedSpare, RemovalSpare)

                LiquidatedList = append(LiquidatedList,
                                        LiquidatedMaterial{
                                                   c.id,
                                                   liquidatedVol,
                                        })


            }
            Removal.RemovalCost = RemovalSpare
        }

    }

    json.NewEncoder(w).Encode(LiquidatedList)

}





func AddOrUpdateContainer(w http.ResponseWriter, r *http.Request) {

    params := mux.Vars(r)
    contId, _ := strconv.Atoi(params["id"])

    var cont Container
    result := json.NewDecoder(r.Body).Decode(&cont)

    if result != nil {
        fmt.Fprintf(os.Stderr, "error=%v\n", result)
        return
    }

    fmt.Printf("container id %d params %v cont.Volume %f cont.HazardPrLitre %f cont.RemovalCostPrLitre %f\n", contId, cont, cont.Volume, cont.HazardPrLitre, cont.RemovalCostPrLitre)

    Containers[contId] = cont

}



