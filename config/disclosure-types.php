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
                        ['pbdf.gemeente.personalData.familyname', 'pbdf.nuts.agb']
                    ]
                ]
            ],
            'valid_authentication' => [
                ['pbdf.nuts.agb', 'pbdf.gemeente.personalData.familyname','pbdf.pbdf.email.email'],
            ],
        ],
        'patient' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        ['pbdf.pbdf.email.email'],
                    ],
                    [
                        ['pbdf.gemeente.personalData.familyname', 'pbdf.gemeente.personalData.initials', 'pbdf.gemeente.personalData.gender', 'pbdf.gemeente.personalData.dateofbirth', 'pbdf.gemeente.personalData.bsn']
                    ]
                ]
            ],
            'valid_authentication' => [
                ['pbdf.gemeente.personalData.familyname', 'pbdf.gemeente.personalData.initials', 'pbdf.gemeente.personalData.gender', 'pbdf.gemeente.personalData.dateofbirth', 'pbdf.gemeente.personalData.bsn'],
            ],
        ],
        'teacher' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        ['pbdf.pbdf.email.email'],
                    ],
                    [
                        ['pbdf.pbdf.surfnet-2.institute', 'pbdf.pbdf.surfnet-2.type', 'pbdf.pbdf.surfnet-2.familyname', 'pbdf.pbdf.surfnet-2.email']
                    ]
                ]
            ],
            'valid_authentication' => [
                ['pbdf.pbdf.surfnet-2.institute', 'pbdf.pbdf.surfnet-2.type', 'pbdf.pbdf.surfnet-2.familyname', 'pbdf.pbdf.surfnet-2.email'],
            ],
        ],
        'student' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        ['pbdf.pbdf.email.email'],
                    ],
                    [
                        ['pbdf.pbdf.surfnet-2.institute', 'pbdf.pbdf.surfnet-2.type', 'pbdf.pbdf.surfnet-2.familyname', 'pbdf.pbdf.surfnet-2.email', 'pbdf.pbdf.surfnet-2.id']
                    ]
                ]
            ],
            'valid_authentication' => [
                ['pbdf.pbdf.surfnet-2.institute', 'pbdf.pbdf.surfnet-2.type', 'pbdf.pbdf.surfnet-2.familyname', 'pbdf.pbdf.surfnet-2.email', 'pbdf.pbdf.surfnet-2.id'],
            ],
        ]


];
