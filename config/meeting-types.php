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
            'irma_disclosure' => 'patient',
            'irma_disclosure_host' => 'medical',
            'attribute_abbreviation' => [
                'pbdf.pbdf.email.email' => '<📧>',
                'pbdf.gemeente.personalData.fullname' => '<BRP> ',
                'pbdf.pbdf.linkedin.firstname' => '<LinkedIn> ',
                'pbdf.pbdf.linkedin.familyname' => ' ',
                'pbdf.pbdf.big.profession' => ' ',
                'pbdf.pbdf.big.specialism' => ' ',
                'pbdf.gemeente.personalData.familyname' => ' ',
                'pbdf.gemeente.personalData.initials' => ' ',
                'pbdf.gemeente.personalData.gender' => ' ',
                'pbdf.gemeente.personalData.dateofbirth' => ' ',
                'pbdf.gemeente.personalData.bsn' => ' '
            ],
            'number_of_participants' => 1

        ],
        'exam' => [
            'irma_disclosure' => 'student',
            'irma_disclosure_host' => 'teacher',
            'attribute_abbreviation' => [
#                'pbdf.pbdf.email.email' => '<📧>',
                'pbdf.pbdf.email.email' => '<Email >',
                'pbdf.gemeente.personalData.fullname' => '<BRP> ',
                'pbdf.pbdf.linkedin.firstname' => '<LinkedIn> ',
                'pbdf.pbdf.linkedin.familyname' => ' ',
                'pbdf.pbdf.big.profession' => ' ',
                'pbdf.pbdf.big.specialism' => ' ',
                'pbdf.pbdf.surfnet-2.institute' => ' ',
                'pbdf.pbdf.surfnet-2.type' => ' ',
                'pbdf.pbdf.surfnet-2.familyname' => ' ',
                'pbdf.pbdf.surfnet-2.id' => ' ',
#                'pbdf.pbdf.surnet-2.email' => '<📧>'
                'pbdf.pbdf.surnet-2.email' => '<Email >'
            ],
            'number_of_participants' => 1

        ]


];
