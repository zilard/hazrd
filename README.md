# Problem statement

Assume you have `n` hazardous liquids `m1, m2, ..., mn` in some inventory.
The hazard classification of each liquid `mi` is given as a positive real number `hi`.
The higher the hazard classification the more dangerous the material is.
The total hazard value of material `mi` at your inventory is given by

    h(mi) = hi * vi

where `vi` is the volume, in liters, of liquid `mi` in your storage.

The insurance cost of your facility are calculated:

          n
    c * MAX  h(mi)
          i=0

Where `c` is a constant

You'd like to reduce your instance cost by removing some of the material from your inventory.
Each liter of liquid `mi` will cost you `ri` to safely remove.
You have `S` USD to spare for the removal and you'd like to spend as much of it as possible
to reduce your insurance cost.

How much of each liquid (in liters) should be removed?


# Delivery

The above problem is solved as a microservice exposing a JSON HTTP API,
with the following endpoints, and the service is packaged as a docker cotainer.

METHOD | QUERY | BODY | DESCRIPTION | RETURN
:--- | :--- | :--- | :--- | :---
**PUT** | `/container/{id}` | `{`&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;<br>&emsp;&emsp;`volume: float,`<br>&emsp;&emsp;`hazardPrLitre: float,`<br>&emsp;&emsp;`removalCostPrLitre: float`<br>`}` | Add or update a&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;<br>container. Where `id` is<br>an integer.
**POST** | `/liquidate` | `{`<br>&emsp;&emsp;`removalCost: float`<br>`}` | Remove liquids from<br>containers such that it<br>satisfies the above<br>problem statement.<br>The `removalCost` parameter<br>corresponds to the `S` USD<br>to spare for removal. | `[`&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;<br>&emsp;&emsp;`{`<br>&emsp;&emsp;&emsp;&emsp;`id: int,`<br>&emsp;&emsp;&emsp;&emsp;`volumeToLiquidate: float`<br>&emsp;&emsp;`{`<br>`]`




# How to build the source code

    make image




# How to run the server

    make run

or

    docker run -p 8080:8080 hazrd-image:latest hazrd


    
# Perform tests

```
curl -i -X PUT -H "Content-Type:application/json" \
http://localhost:8080/container/1 -d '{"volume": 1, "hazardPrLitre": 2, "removalCostPrLitre": 3}'
```

```
curl -i -X PUT -H "Content-Type:application/json" \
http://localhost:8080/container/2 -d '{"volume": 11, "hazardPrLitre": 22, "removalCostPrLitre": 33}'
```

```
curl -i -X GET -H "Content-Type:application/json" \
http://localhost:8080/showcontainers

{"1":{"volume":1,"hazardPrLitre":2,"removalCostPrLitre":3},"2":{"volume":11,"hazardPrLitre":22,"removalCostPrLitre":33}}
```

```
curl -i -X POST -H "Content-Type:application/json" \
http://localhost:8080/liquidate -d '{"removalCost": 3}'

[{"id":1,"volumeToLiquidate":1}]
```

```
curl -i -X GET -H "Content-Type:application/json" \
http://localhost:8080/showcontainers

{"2":{"volume":11,"hazardPrLitre":22,"removalCostPrLitre":33}}
```

```
curl -i -X PUT -H "Content-Type:application/json" \
http://localhost:8080/container/3 -d '{"volume": 2, "hazardPrLitre": 3, "removalCostPrLitre": 5}'
```

```
curl -i -X GET -H "Content-Type:application/json" \
http://localhost:8080/showcontainers

{"2":{"volume":11,"hazardPrLitre":22,"removalCostPrLitre":33},"3":{"volume":2,"hazardPrLitre":3,"removalCostPrLitre":5}}
```

```
curl -i -X POST -H "Content-Type:application/json" \
http://localhost:8080/liquidate -d '{"removalCost": 38}'

[{"id":2,"volumeToLiquidate":1.1515151}]
```

```
curl -i -X GET -H "Content-Type:application/json" \
http://localhost:8080/showcontainers

{"2":{"volume":9.848485,"hazardPrLitre":22,"removalCostPrLitre":33},"3":{"volume":2,"hazardPrLitre":3,"removalCostPrLitre":5}}
```

```
curl -i -X POST -H "Content-Type:application/json" \
http://localhost:8080/liquidate -d '{"removalCost": 380}'

[{"id":2,"volumeToLiquidate":9.848485},{"id":3,"volumeToLiquidate":2}]
```

```
curl -i -X GET -H "Content-Type:application/json" \
http://localhost:8080/showspare

{"removalCost":45}
```



