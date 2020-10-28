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
                        ['pbdf.sidn-pbdf.email.email'],
                    ],
                    [
                        [
                        'pbdf.gemeente.personalData.fullname'
                        ],
                                    [
                        'pbdf.pbdf.linkedin.firstname',
                        'pbdf.pbdf.linkedin.familyname'
                        ]
                    ]
                ]
            ],
            'valid_authentication' => [
                [
                    'pbdf.gemeente.personalData.fullname',
                    'pbdf.pbdf.email.email'
                ],
                [
                    'pbdf.pbdf.linkedin.firstname',
                    'pbdf.pbdf.linkedin.familyname',
                    'pbdf.pbdf.email.email'
                ],
                [
                    'pbdf.gemeente.personalData.fullname',
                    'pbdf.sidn-pbdf.email.email'
                ],
                [
                    'pbdf.pbdf.linkedin.firstname',
                    'pbdf.pbdf.linkedin.familyname',
                    'pbdf.sidn-pbdf.email.email'
                ]
            ],
            'name' => [
                [
                    'pbdf.gemeente.personalData.fullname'
                ],
                [
                    'pbdf.pbdf.linkedin.firstname',
                    'pbdf.pbdf.linkedin.familyname'
                ]
            ],
            'email' => ['pbdf.sidn-pbdf.email.email', 'pbdf.pbdf.email.email'],
            'hoster_message' => 'Join as host',
            'client_message' => 'Join as participant'
        ],
        'medical' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        [
            'pbdf.pbdf.email.email'
            ],
                    ],
                    [
                        [
                            'pbdf.nuts.agb.agbcode',
                //            'pbdf.pbdf.mobilenumber.mobilenumber',
                            'pbdf.gemeente.personalData.familyname'
                        ]
                    ]
                ]
            ],
            'valid_authentication' => [
                [
                    'pbdf.nuts.agb.agbcode',
            //        'pbdf.pbdf.mobilenumber.mobilenumber',
                    'pbdf.gemeente.personalData.familyname',
                    'pbdf.pbdf.email.email'
                ],
                [
                    'pbdf.nuts.agb.agbcode',
            //        'pbdf.pbdf.mobilenumber.mobilenumber',
                    'pbdf.gemeente.personalData.familyname',
                    'pbdf.sidn-pbdf.email.email'
                ],
            ],
            'name' => [
                ['pbdf.gemeente.personalData.familyname'],
            ],
            'email' => ['pbdf.sidn-pbdf.email.email', 'pbdf.pbdf.email.email'],
            'hoster_message' => 'Join as doctor',
            'client_message' => 'Join as doctor'
        ],
        'patient' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        ['pbdf.pbdf.email.email'],
                        ['pbdf.sidn-pbdf.email.email']
                    ],
                    [
                        [
                      'pbdf.gemeente.personalData.initials',
                      'pbdf.gemeente.personalData.gender',
                      'pbdf.gemeente.personalData.dateofbirth',
                      'pbdf.gemeente.personalData.familyname'
                    ]
                    ]
                ]
            ],
            'valid_authentication' => [
                [
                    'pbdf.gemeente.personalData.familyname',
                    'pbdf.gemeente.personalData.initials',
                    'pbdf.gemeente.personalData.gender',
                    'pbdf.gemeente.personalData.dateofbirth'
                ],
            ],
            'name' => [
                [
                    'pbdf.gemeente.personalData.initials',
                    'pbdf.gemeente.personalData.familyname'
                ],
            ],
            'email' => ['pbdf.sidn-pbdf.email.email', 'pbdf.pbdf.email.email'],
            'hoster_message' => 'Join as patient',
            'client_message' => 'Join as patient'
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
                [
                    'pbdf.pbdf.surfnet-2.institute',
                    'pbdf.pbdf.surfnet-2.type',
                    'pbdf.pbdf.surfnet-2.familyname',
                    'pbdf.pbdf.surfnet-2.email'],
                ],
            'name' => ['pbdf.pbdf.surfnet-2.familyname'],
            'email' => ['pbdf.pbdf.surfnet-2.email'],
            'hoster_message' => 'Join as teacher',
            'client_message' => 'Join as teacher'
        ],
        'student' => [
            'irma_disclosure' => [
                '@context' => 'https://irma.app/ld/request/disclosure/v2',
                'disclose' => [
                    [
                        [
                            ['type' => 'pbdf.pbdf.surfnet-2.institute', 'value' => null],
                            ['type' => 'pbdf.pbdf.surfnet-2.type', 'value' => 'student'],
                            ['type' => 'pbdf.pbdf.surfnet-2.familyname', 'value' => null],
                            ['type' => 'pbdf.pbdf.surfnet-2.email', 'value' => null]
                        ]
                    ]
                ]
            ],
            'valid_authentication' => [
                [
                    'pbdf.pbdf.surfnet-2.institute',
                    'pbdf.pbdf.surfnet-2.type',
                    'pbdf.pbdf.surfnet-2.familyname',
                    'pbdf.pbdf.surfnet-2.email'
                ],
            ],
            'name' => [['pbdf.pbdf.surfnet-2.familyname']],
            'email' => 'pbdf.pbdf.surfnet-2.email',
            'hoster_message' => 'Join as student',
            'client_message' => 'Join as student'
        ]


];
