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
            'irma_disclosure' => 'default',
            'attribute_abbreviation' => [
                'pbdf.pbdf.email.email' => '[📧]',
                'pbdf.gemeente.personalData.fullname' => '[BRP] ',
                'pbdf.pbdf.linkedin.firstname' => '[LinkedIn] ',
                'pbdf.pbdf.linkedin.familyname' => ' '
            ]

        ],
        'medical_consult' => [
            'irma_disclosure' => 'default',
            'irma_disclosure_host' => 'medical',
            'attribute_abbreviation' => [
                'pbdf.pbdf.email.email' => '<📧>',
                'pbdf.gemeente.personalData.fullname' => '<BRP> ',
                'pbdf.pbdf.linkedin.firstname' => '<LinkedIn> ',
                'pbdf.pbdf.linkedin.familyname' => ' ',
                'pbdf.pbdf.big.profession' => ' ',
                'pbdf.pbdf.big.specialism' => ' '
            ],
            'number_of_participants' => 1

        ],
        'exam' => [
            'irma_disclosure' => 'student',
            'irma_disclosure_host' => 'teacher',
            'attribute_abbreviation' => [
                'pbdf.pbdf.email.email' => '<📧>',
                'pbdf.gemeente.personalData.fullname' => '<BRP> ',
                'pbdf.pbdf.linkedin.firstname' => '<LinkedIn> ',
                'pbdf.pbdf.linkedin.familyname' => ' ',
                'pbdf.pbdf.big.profession' => ' ',
                'pbdf.pbdf.big.specialism' => ' '
            ],
            'number_of_participants' => 1

        ]


];
