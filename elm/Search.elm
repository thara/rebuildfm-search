module Search exposing (..)

import Date
import Debug
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)
import Html.Lazy exposing (lazy)
import Task
import Http
import Json.Decode as Json


-- MAIN


main : Program Flags Model Msg
main =
    Html.programWithFlags
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }


-- MODEL
type alias Model =
  { queryCastName: String
  , queryTitle: String
  , queryNote: String
  , episodes : List Episode
  , apiBaseUrl : String
  }


type alias Cast =
    { name: String
    , uri: String}


type alias Flags =
    { apiBaseUrl : String
    }


type alias Episode =
  {
    no: Int
  , title: String
  , date_published: String
  , episodeUrl: String
  , description: String
  , starring: List Cast
  }


init : Flags -> ( Model, Cmd Msg )
init flags =
    (
      { queryCastName = ""
      , queryTitle = ""
      , queryNote = ""
      , episodes = [
--          [ { no = 155
--            , title = "155: I Am Your Grandfather (Matz)"
--            , date_published = "Aug 21 2016"
--            , episodeUrl = "http://rebuild.fm/155/"
--            , description = "まつもとゆきひろさんをゲストに迎えて、RubyKaigi, Ruby 2.4, プロダクティビティ、mruby, Elixir, 自作言語などについて話しました。"
--            , starring = ["matz", "miyagawa"]
--            }
--          , { no = 154
--            , title = "Aftershow 154: Sick Of Experiments (kazuho)"
--            , date_published = "Aug 18 2016"
--            , episodeUrl = "http://rebuild.fm/154/"
--            , description = "Kazuho Oku さんと、CPU実験、大学入試などについて話しました。"
--            , starring = ["kazuho", "miyagawa"]
--            }
          ]
      , apiBaseUrl = flags.apiBaseUrl
      },
      Cmd.none
    )



-- MESSAGES
type Msg =
  InputCastName String
  | InputTitle String
  | InputNote String
  | SearchEpisodes
  | ReceiveEpisodes (Result Http.Error (List Episode))



-- VIEW


view : Model -> Html Msg
view model =
  body
    []
    [
      div
        [ class "container"]
        [ div
            [ class "row jumbotron" ]
            [ h1
                [ class "display-3" ]
                [ text "Rebuild.fm Search" ]
              , viewInput model
            ]
        , div
            [ class "row" ]
            [ lazy viewEpisodes model.episodes ]

        ]
    , footer
        [ class "footer" ]
        [ div
            [ class "container" ]
            [ span [ class "text-muted"] [ text "Rebuild.fm Search © 2016" ] ]
        ]
    ]


viewInput : Model -> Html Msg
viewInput model =
  div
    []
    [ Html.form
        [ class "form-inline" ]
        [ div
            [ class "form-group" ]
            [ label
                [ class "sr-only", for "input-cast-name" ]
                [ text "Cast" ]
             , input
                 [ type_ "text"
                 , id "input-cast-name"
                 , class "form-control"
                 , placeholder "Cast Name"
                 , onInput InputCastName ] []
            ]
        , div
             [ class "form-group" ]
             [ label
                 [ class "sr-only", for "input-title" ]
                 [ text "Title" ]
              , input
                  [ type_ "text"
                  , id "input-title"
                  , class "form-control"
                  , placeholder "Title"
                  , onInput InputTitle ] []
             ]
        , div
             [ class "form-group" ]
             [ label
                 [ class "sr-only", for "input-note" ]
                 [ text "Note" ]
              , input
                  [ type_ "text"
                  , id "input-note"
                  , class "form-control"
                  , placeholder "Note"
                  , onInput InputNote ] []
             ]
        , button
            [ type_ "button"
            , class "btn btn-primary"
            , onClick SearchEpisodes ]
            [ text "Search" ]
        ]
    ]


viewEpisodes : List Episode -> Html Msg
viewEpisodes episodes =
  div
    [] <|
    List.map viewEpisode episodes


viewEpisode : Episode -> Html Msg
viewEpisode episode =
  div
    [ class "card card-block" ]
    [ h4
        [ class "card-title" ]
        [ a [ href episode.episodeUrl ] [ text episode.title ] ]
    , p
        [ class "card-text" ]
        [ small [ class "text-muted"] [ text episode.date_published ] ]
    , p
        [ class "card-text" ]
        [ text episode.description ]
    , p
        [ class "card-text" ]
        [ h6 [] [ text "Starring" ]
        , ul [] <|
            List.map viewEpisodeStarring episode.starring
        ]
    ]


viewEpisodeStarring : Cast -> Html Msg
viewEpisodeStarring cast =
    li [] [
      a [ href cast.uri ] [ text cast.name ]
    ]


-- UPDATE

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  case msg of
    InputCastName s -> ( {model | queryCastName = s}, Cmd.none)
    InputTitle s -> ( {model | queryTitle = s}, Cmd.none)
    InputNote s -> ( {model | queryNote = s}, Cmd.none)

    SearchEpisodes ->
        (model, searchEpisodes model)
    ReceiveEpisodes (Ok episodes) ->
        ( { model | episodes = episodes}, Cmd.none)
    ReceiveEpisodes (Err _) ->
        (model, Cmd.none)



searchEpisodes : Model -> Cmd Msg
searchEpisodes model =
    let
        url = model.apiBaseUrl ++ "episodes?cast_name=" ++ model.queryCastName ++ "&title=" ++ model.queryTitle ++ "&note=" ++ model.queryNote
        request =
            Http.request
                { method = "GET"
                , headers = []
                , url = url
                , body = Http.emptyBody
                , expect = Http.expectJson decodeEpisodes
                , timeout = Nothing
                , withCredentials = False
                }
    in
        Http.send ReceiveEpisodes request


decodeEpisodes : Json.Decoder (List Episode)
decodeEpisodes =
    Json.field "episodes" (Json.list decodeEpisode)


decodeEpisode : Json.Decoder Episode
decodeEpisode =
    Json.map6 Episode
        (Json.succeed 42)
        (Json.field "title" Json.string)
        (Json.succeed "2016-01-01")
        (Json.field "link" Json.string)
        (Json.field "subtitle" Json.string)
        (Json.field "casts" (Json.list decodeCast))


decodeCast : Json.Decoder Cast
decodeCast =
    Json.map2 Cast
        (Json.field "name" Json.string)
        (Json.field "uri" Json.string)


-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none
