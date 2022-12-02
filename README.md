# BT_API
It's api service which save dayly picture and description from APOD.
API have 3 endpoints

## API
|Method         | Path           | Operation  |
| :-----------: | :------------: | :--------: | 
| GET           | /articles     | data:{[{"date":string, "title":string,"explanation":string,"url":string}]}|
| GET           | /article/{data string}  | data:{"date":string, "title":string,"explanation":string,"url":string}|
| GET           | /picture/{data string}  | image/jpeg |


# Get start
1. git clone https://github.com/GritselMaks/BT_API
2. docker-compose up

# Note
If you want to add pictures from the previous month, you should uncommitted func s.AddContent().
New picture added every day
