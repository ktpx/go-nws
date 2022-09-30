complete -e gonws
complete -c gonws -f
complete -c gonws -o area -x -d 'Area (AR,AH,CA,FL...)'
complete -c gonws -s x -x -a "alerts count" -d 'Report type'
complete -c gonws -s r -x -d 'Region'
complete -c gonws -s c -x -a "Observed Likely Possible Unlikely" -d 'Certainty'
complete -c gonws -s r -x -a "AL AT GL PA PI" -d 'Marine Region Code'
complete -c gonws -o rt -x -a "land marine" -d 'Marine Region Type'
complete -c gonws -s s -x -a "actual excercise system test draft" -d 'Status'
complete -c gonws -s z -x -d 'Zone Code'
complete -c gonws -s e -x -d 'Event Name'
complete -c gonws -s u -x -a "Immediate Expected Future Past Uknown" -d 'Urgency'
