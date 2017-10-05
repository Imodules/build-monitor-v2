module Api exposing (..)

import Decoders exposing (buildTypesDecoder, projectsDecoder)
import Http
import Json.Decode as Decode exposing (Decoder)
import Msgs exposing (Msg)
import RemoteData
import Types exposing (..)
import Urls


authHeader : Token -> Http.Header
authHeader token =
    Http.header "Authorization" ("Bearer " ++ token)


authGet : String -> Token -> Decode.Decoder a -> Http.Request a
authGet url token decoder =
    Http.request
        { method = "GET"
        , headers = [ authHeader token ]
        , url = url
        , body = Http.emptyBody
        , expect = Http.expectJson decoder
        , timeout = Nothing
        , withCredentials = False
        }


delete : String -> Decode.Decoder a -> Http.Request a
delete url decoder =
    Http.request
        { method = "DELETE"
        , headers = []
        , url = url
        , body = Http.emptyBody
        , expect = Http.expectJson decoder
        , timeout = Nothing
        , withCredentials = False
        }


post : String -> Http.Body -> Decode.Decoder a -> Http.Request a
post url body decoder =
    Http.request
        { method = "POST"
        , headers = []
        , url = url
        , body = body
        , expect = Http.expectJson decoder
        , timeout = Nothing
        , withCredentials = False
        }


authPost : String -> Token -> Http.Body -> Decode.Decoder a -> Http.Request a
authPost url token body decoder =
    Http.request
        { method = "POST"
        , headers = [ authHeader token ]
        , url = url
        , body = body
        , expect = Http.expectJson decoder
        , timeout = Nothing
        , withCredentials = False
        }


authEmptyPost : String -> Token -> Http.Request ()
authEmptyPost url token =
    Http.request
        { method = "POST"
        , headers = [ authHeader token ]
        , url = url
        , body = Http.emptyBody
        , expect = Http.expectStringResponse (\_ -> Ok ())
        , timeout = Nothing
        , withCredentials = False
        }


authPut : String -> Token -> Http.Body -> Decode.Decoder a -> Http.Request a
authPut url token body decoder =
    Http.request
        { method = "PUT"
        , headers = [ authHeader token ]
        , url = url
        , body = body
        , expect = Http.expectJson decoder
        , timeout = Nothing
        , withCredentials = False
        }


patch : String -> Http.Body -> Http.Request ()
patch url body =
    Http.request
        { method = "PATCH"
        , headers = []
        , url = url
        , body = body
        , expect = Http.expectStringResponse (\_ -> Ok ())
        , timeout = Nothing
        , withCredentials = False
        }


fetchProjects : String -> Cmd Msg
fetchProjects baseApiUrl =
    Http.get (Urls.projects baseApiUrl) projectsDecoder
        |> RemoteData.sendRequest
        |> Cmd.map Msgs.OnFetchProjects


fetchBuildTypes : String -> Cmd Msg
fetchBuildTypes baseApiUrl =
    Http.get (Urls.buildTypes baseApiUrl) buildTypesDecoder
        |> RemoteData.sendRequest
        |> Cmd.map Msgs.OnFetchBuildTypes


refreshServerProjects : String -> Token -> Cmd Msg
refreshServerProjects baseApiUrl token =
    let
        request =
            authEmptyPost (Urls.refreshServerProjects baseApiUrl) token
    in
        Http.send Msgs.OnRefreshServerProjects request
