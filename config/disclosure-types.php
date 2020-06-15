<?php

return [

    /*
    |--------------------------------------------------------------------------
    | disclosure-types
    |--------------------------------------------------------------------------
    |
    | This file is for storing the disclosure types for this IRMA meet project
    |
    */
        'default' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        ['pbdf.pbdf.email.email'],
                    ],
                    [
                        ['pbdf.gemeente.personalData.fullname'],
                        ['pbdf.pbdf.linkedin.firstname', 'pbdf.pbdf.linkedin.familyname']
                    ]
                ]
            ],
            'valid_authentication' => [
                ['pbdf.gemeente.personalData.fullname','pbdf.pbdf.email.email'],
                ['pbdf.pbdf.linkedin.firstname', 'pbdf.pbdf.linkedin.familyname', 'pbdf.pbdf.email.email']
            ],
        ],
        'medical' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        ['pbdf.pbdf.email.email'],
                    ],
                    [
                        ['pbdf.gemeente.personalData.familyname', 'pbdf.nuts']
                    ]
                ]
            ],
            'valid_authentication' => [
                ['pbdf.nuts', 'pbdf.gemeente.personalData.familyname','pbdf.pbdf.email.email'],
            ],
        ]


];
