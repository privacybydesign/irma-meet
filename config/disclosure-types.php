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
            'name' => [
                ['pbdf.gemeente.personalData.fullname'],
                ['pbdf.pbdf.linkedin.firstname', 'pbdf.pbdf.linkedin.familyname']
            ],
            'email' => 'pbdf.pbdf.email.email'
        ],
        'medical' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        ['pbdf.pbdf.email.email'],
                    ],
                    [
                        [
                          'pbdf.gemeente.personalData.familyname',
            #             'pbdf.nuts.agb.agbcode',
                          'pbdf.gemeente.personalData.bsn'
                        ]
                    ]
                ]
            ],
            'valid_authentication' => [
                #['pbdf.nuts.agb.agbcode', 'pbdf.gemeente.personalData.familyname', 'pbdf.pbdf.email.email'],
                ['pbdf.gemeente.personalData.familyname', 'pbdf.pbdf.email.email'],
            ],
            'name' => [
                ['pbdf.gemeente.personalData.familyname'],
            ],
            'email' => 'pbdf.pbdf.email.email'
        ],
        'patient' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        ['pbdf.pbdf.email.email'],
                    ],
                    [
                        [
              'pbdf.gemeente.personalData.familyname',
              'pbdf.gemeente.personalData.initials',
              'pbdf.gemeente.personalData.gender',
              'pbdf.gemeente.personalData.dateofbirth',
              'pbdf.gemeente.personalData.bsn'
            ]
                    ]
                ]
            ],
            'valid_authentication' => [
                ['pbdf.gemeente.personalData.familyname', 'pbdf.gemeente.personalData.initials', 'pbdf.gemeente.personalData.gender', 'pbdf.gemeente.personalData.dateofbirth', 'pbdf.gemeente.personalData.bsn'],
            ],
            'name' => [
                ['pbdf.gemeente.personalData.initials', 'pbdf.gemeente.personalData.familyname'],
            ],
            'email' => 'pbdf.pbdf.email.email'
        ],
        'teacher' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        [
              ['type' => 'pbdf.pbdf.surfnet-2.institute', 'value' => null],
              ['type' => 'pbdf.pbdf.surfnet-2.type', 'value' => 'employee'],
              ['type' => 'pbdf.pbdf.surfnet-2.familyname', 'value' => null],
              ['type' => 'pbdf.pbdf.surfnet-2.email', 'value' => null]
            ]
                    ]
                ]
            ],
            'valid_authentication' => [
                ['pbdf.pbdf.surfnet-2.institute', 'pbdf.pbdf.surfnet-2.type', 'pbdf.pbdf.surfnet-2.familyname', 'pbdf.pbdf.surfnet-2.email'],
            ],
            'name' => [['pbdf.pbdf.surfnet-2.familyname']],
            'email' => 'pbdf.pbdf.surfnet-2.email'
        ],
        'student' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        [
              ['type' => 'pbdf.pbdf.surfnet-2.institute', 'value' => null],
              ['type' => 'pbdf.pbdf.surfnet-2.type', 'value' => 'employee'],
              ['type' => 'pbdf.pbdf.surfnet-2.familyname', 'value' => null],
              ['type' => 'pbdf.pbdf.surfnet-2.email', 'value' => null]
            ]
                    ]
                ]
            ],
            'valid_authentication' => [
                ['pbdf.pbdf.surfnet-2.institute', 'pbdf.pbdf.surfnet-2.type', 'pbdf.pbdf.surfnet-2.familyname', 'pbdf.pbdf.surfnet-2.email'],
            ],
            'name' => [['pbdf.pbdf.surfnet-2.familyname']],
            'email' => 'pbdf.pbdf.surfnet-2.email'
        ]


];
