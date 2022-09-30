# go-nws

Go CLI Client for NWS (National Weather Service) Alerts

The client queries and lists alerts from the NWS API (active alerts) endpoint.  

### Example

```
$ gonws -x alerts -area FL -s Extreme -c Observed,Likely

$ gonws -x count

```

Note that that argument attributes are case sensitive (ie, `Extreme`). Check
`-h`for more options.

### Example report for Alerts

```
==================================================
Event    : Flash Flood Watch
Headline : Flash Flood Watch issued September 29 at 11:57AM PDT until September 30 at 3:00PM PDT by NWS Spokane WA
Category : Met
Msgtype  : Alert
Urgency  : Future
Certainty: Possible
Type     : wx:Alert
Sent     : 2022-09-29T11:57:00-07:00
Effective: 2022-09-29T11:57:00-07:00
Onset    : 2022-09-29T11:57:00-07:00
Expires  : 2022-09-29T15:00:00-07:00
Sender   : NWS Spokane WA (w-nws.webmaster@noaa.gov)
AreaDesc : Pend Oreille, WA
Description  :
* WHAT...Flooding caused by upstream dam releases of Box Canyon Dam
is possible on the Pend Oreille River.

* WHERE...A portion of Northeast Washington, including Pend Oreille
county.

* WHEN...Through this afternoon.

* IMPACTS...Increased releases may result in flash flooding of
low-lying areas below the dam, especially to Metaline and Metaline
Falls.

* ADDITIONAL DETAILS...
- Flows have increased from the Box Canyon Dam from 16 kcfs to
60 kcfs.
- http://www.weather.gov/safety/flood
--------------------------------------------------
Instructions : You should monitor later forecasts and be prepared to take action
should Flash Flood Warnings be issued.
==================================================
```
### Count Report

```
----------------------------------------
Total Alert Count Report
----------------------------------------
Total  : 495
Land   : 490
Marine : 178
----------------------------------------
Alerts per Area
----------------------------------------
 AK :   7  AL :   2  AM :  14  AN :  33
 CO :   1  DE :   3  FL :  48  GA :  41
 GM :   5  IA :   1  KS :   1  LE :   4
 MD :   8  MT :   1  NC : 124  NE :   1
 NJ :   4  OR :   1  PK : 111  PZ :  11
 SC :  71  SD :   1  TN :   1  TX :   7
 VA :  30  VI :   1  WA :   1  WI :   1
 WV :   2  WY :   1
```
