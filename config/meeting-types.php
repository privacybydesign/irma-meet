<?php

return [

    /*
    |--------------------------------------------------------------------------
    | meetintTypes
    |--------------------------------------------------------------------------
    |
    | This file is for storing the meeting types for this IRMA meet project
    |
    */
        'free' => [
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
                ],
            ],
            'valid_authentication' => [
                ['pbdf.pbdf.email.email', 'pbdf.gemeente.personalData.fullname'],
                ['pbdf.pbdf.email.email', 'pbdf.pbdf.linkedin.firstname', 'pbdf.pbdf.linkedin.familyname']
            ],
            'attribute_abbreviation' => [
                'pbdf.pbdf.email.email' => '@',
                'pbdf.gemeente.personalData.fullname' => '[BRP]',
                'pbdf.pbdf.linkedin.firstname' => 'linkedIn',
                'pbdf.pbdf.linkedin.familyname' => ' '
            ]

        ],
        'ggz' => [

        ]


];
