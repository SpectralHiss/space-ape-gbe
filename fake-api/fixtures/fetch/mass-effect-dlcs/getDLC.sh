DLCs="3100 3101 3102 3184 3185 3186 3288 3289 3290 3431 3432 3433 3561 3562 3563 3706 3742 3866 3867 3868 3869 3870 3871 3919 3920 3921 3922 3929"


for DLC in $DLCs; do

 wget --header="Accept: text/html" --user-agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:21.0) Gecko/20100101 Firefox/21.0" \
   "https://giantbomb.com/api/dlc/${DLC}?api_key=${API_KEY}&format=json" -O $DLC.json
done
