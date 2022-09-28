# go-nws

Go CLI Client for NWS (National Weather Service) Alerts

The client queries and lists alerts from the NWS API (active alerts) endpoint.  

##Example

```
$ ./gonws -area FL -s Extreme -c Observed,Likely
```

All options, are just that. Optional.  No arguments
will list all alerts. Not that argument fields are 
case sensitive (ie, `Extreme`) as per the help file.

## Example output


```
====================================================================
Event    : Storm Surge Warning
Headline : Storm Surge Warning issued September 28 at 11:03AM EDT by
Category : Met
Urgency  : Immediate
Type     : wx:Alert
Sent     : 2022-09-28T11:03:00-04:00
Effective: 2022-09-28T11:03:00-04:00
Onset    : 2022-09-28T11:03:00-04:00
Expires  : 2022-09-28T19:15:00-04:00
Sender   : NWS Miami FL (w-nws.webmaster@noaa.gov)
Msgtype  : Update
Desc     :
* LOCATIONS AFFECTED
- Flamingo
- Cape Sable
- Loop Road

* WIND
- LATEST LOCAL FORECAST: Equivalent Tropical Storm force wind
- Peak Wind Forecast: 30-40 mph with gusts to 75 mph
- Window for Tropical Storm force winds: through the next few
hours

- THREAT TO LIFE AND PROPERTY THAT INCLUDES TYPICAL FORECAST
UNCERTAINTY IN TRACK, SIZE AND INTENSITY: Potential for wind 39
to 57 mph
- The wind threat has decreased from the previous assessment.
- PLAN: Plan for hazardous wind of equivalent tropical storm
force.
- PREPARE: Last minute efforts to protect property should now
be complete. The area remains subject to limited wind
damage.
- ACT: Now is the time to shelter from hazardous wind.

```
