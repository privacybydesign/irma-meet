<?php

namespace App\Http\Middleware;

use Closure;
use Session;

class IrmaAuthenticate
{
    /**
     * Handle an incoming request.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  \Closure  $next
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        $token = Session::get('irmasession', '');
        if ( $request->route()->getName() === 'irma_auth.start' ) {
            //start session
            return response($this->start($request));
        } else if ( $token === '' ) {
            return redirect(route('home'));
        } else {
            //verify
            $result = $this->irma_get_session_result($token);
            if ( $result && $result->proofStatus == 'VALID' ) {
                $test = $result->disclosed[0][0]->rawvalue;
                $validated_email = $result->disclosed[0][0]->rawvalue;
                session(['validated_email' => $validated_email]);
            } else {
                //TODO present error
                return redirect(route('home'));
            }
        }
        return $next($request);
    }


    private function start($request)
    {
        $irmasession = $this->irma_start_session([
            '@context' => 'https://irma.app/ld/request/disclosure/v2',
            'disclose' => [
                [
                    ['pbdf.pbdf.email.email'],
                ],
            ],
        ]);

        session(['irmasession' => $irmasession->token]);
        
        return json_encode($irmasession->sessionPtr);
    }

    private function irma_server_call($method, $suburl, $payload=NULL) {
        if (! isset($payload)) {
            $api_call = array(
                'http' => array(
                    'method' => $method,
                    //'header' => "Authorization: " . IRMA_API_TOKEN . "\r\n",
                )
            );
        } else {
            $json_payload = json_encode($payload);
            
            $api_call = array(
                'http' => array(
                    'method' => $method,
                    'header' => "Content-Type: application/json\r\n"
                        . "Content-Length: " . strlen($json_payload) . "\r\n",
                        //. "Authorization: " . IRMA_API_TOKEN ."\r\n",
                    'content' => $json_payload
                )
            );
        }
        
        $url = env("IRMA_SERVER_URL") . $suburl;
        $file_headers = @get_headers($url);
        if(!$file_headers || $file_headers[0] == 'HTTP/1.0 404 Not Found' || $file_headers[0] == 'HTTP/1.0 400 Bad Request') {
            $response_data = false;
        } else {
            $response = file_get_contents($url, false, stream_context_create($api_call));
            if (! $response) {
                // TODO: integrate with whatever we use for logging
                http_response_code(501);
                echo 'Internal server error';
                exit();
            }

            $response_data = json_decode($response);
            if ($response_data === NULL) {
                // TODO: logging
                http_response_code(502);
                echo 'Internal server error';
                exit();
            }
        }


        return $response_data;
    }

    private function irma_start_session($request) {
        return $this->irma_server_call('POST', '/session', $request);
    }

    private function irma_get_session_result($token) {
        return $this->irma_server_call('GET', '/session/' . $token . '/result');
    }

}
