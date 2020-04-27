<?php

namespace App\Http\Middleware;


class IrmaAuthenticate
{

    public function handle($request, Closure $next)
    {
        // Perform action

        return $next($request);
    }

    private function start()
    {
        $irmasession = $this->irma_start_session([
            '@context' => 'https://irma.app/ld/request/disclosure/v2',
            'disclose' => [
                [
                    ['pbdf.pbdf.email.email'],
                ],
            ],
        ]);

        $_SESSION['irmasession'] = $irmasession->token;

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
        
        $response = file_get_contents(env("IRMA_SERVER_URL") . $suburl, false, stream_context_create($api_call));
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

        return $response_data;
    }

    private function irma_start_session($request) {
        return $this->irma_server_call('POST', '/session', $request);
    }

    private function irma_get_session_result($token) {
        return $this->irma_server_call('GET', '/session/' . $token . '/result');
    }

}
