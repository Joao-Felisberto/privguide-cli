import json

import yaml

from test_case_generator.dfd import DataType, ExternalEntity, Process, DataStore, DataStored, DataFlow, DFD
from test_case_generator.dpia import DPO, PersonalDatum, DPIA, Risk, Processing, Purpose, SupervisoryAuthorityVeredict


def to_yaml(data, fname):
    data = json.dumps(data, default=lambda x: {
        k.replace("_", " ").rstrip(): x.__dict__[k]
        for k in x.__dict__
    })

    with open(fname, 'w') as f:
        f.write(yaml.dump(json.loads(data)))


if __name__ == '__main__':
    dt1 = DataType(
        "message",
        [],
        "eternal",
        [
            "dpia:confidential",
            "dpia:personal",
        ],
    )
    dt2 = DataType(
        "AccountId",
        [":AccountId"],
        "eternal",
        ["dpia:personal"],
    )
    dt3 = DataType(
        "Account",
        [":AccountId"],
        "eternal",
        [
            "dpia:confidential",
            "dpia:personal",
        ],
    )

    ee1 = ExternalEntity(
        "dpia:User",
        [":message"],
        [":message"],
        ["Portugal"],
        [],
        ["dpia:human"],
        ">16",
        False,
        [],
        [],
    )
    ee2 = ExternalEntity(
        "dpia:User",
        [":message"],
        [":message"],
        ["Portugal"],
        [],
        ["dpia:external system"],
        None,
        False,
        [],
        [],
    )

    proc1 = Process(
        "send message",
        [":message"],
        [":message"],
        ["Portugal"],
        [],
        ["dpia:message routing"],
        [],
        [],
    )

    data1 = DataStored(
        ":message",
        "eternal",
        ":C store message",
        ":R store message",
        ":U store message",
        ":D store message",
    )

    ds1 = DataStore(
        "message db",
        [data1],
        ["Portugal"],
        [],
        [],
        []
    )

    df1 = DataFlow(
        "C message",
        "dpia:User",
        ":send message",
        [":message"],
        "signal",
        "1m",
        1,
        [],
        []
    )
    df2 = df1.clone(id_="R message")
    df3 = df1.clone(id_="U message")
    df4 = df1.clone(id_="D message")

    df5 = DataFlow(
        "C store message",
        ":send message",
        ":message db",
        [":message"],
        "signal",
        "1m",
        1,
        [],
        []
    )
    df6 = df5.clone(id_="R store message")
    df7 = df5.clone(id_="U store message")
    df8 = df5.clone(id_="D store message")

    dfd = DFD(
        [dt1, dt2, dt3],
        [ee1, ee2],
        [proc1],
        [ds1],
        [df1, df2, df3, df4, df5, df6, df7, df8]
    )

    dpo1 = DPO(
        ":The",
        "the@email.com"
    )
    dpo2 = DPO(
        ":Man",
        "manemail.com"
    )

    pd1 = PersonalDatum(
        "dfd:message",
        ":personal",
        [],
        [],
        [":User"],
        "2d",
        False,
        [":message routing"],
        []
    )
    pd2 = PersonalDatum(
        "dfd:AccountId",
        ":personal",
        [],
        [],
        [":User"],
        "2d",
        False,
        [],
        []
    )
    pd3 = pd2.clone(id_="dfd:Account")

    risk1 = Risk(
        "Risk 1",
        0,
        1,
        []
    )

    purpose1 = Purpose(
        ":message routing",
        True,
        True,
        True,
    )

    processing1 = Processing(
        id_="dfd:send message",
        requires_new_technologies=False,
        risk_to_rights_and_freedoms_of=[":User"],
        required_for_contract=[],
        legally_mandated=False,
        vital_interest=[],
        public_interest=False,
        is_official_authority=False,
        legitimate_interest=[":User"],
        professional_secrecy=False,
        scores_users=True,
        automated_decisions=True,
        legal_impact_for_the_user=True,
        systematic_monitoring=True,
        large_scale_processing=False,
        lawful=True,
        fair=True,
        transparent=True,
        specific=True,
        explicit=True,
        legitimate=True,
        purposes=[purpose1],
        risks=[f":{risk1.id_}"],
        supervisory_authority_veredict=SupervisoryAuthorityVeredict(
            [":Supervisor"],
            True
        )
    )

    dpia = DPIA(
        ":last update",
        [":Someone", ":Else"],
        [dpo1, dpo2],
        [pd1, pd2, pd3],
        [risk1],
        [":message routing"],
        [],
        [processing1]
    )

    to_yaml(dfd, "../.devprivops/tests/a/out.dfd.yml")
    to_yaml(dpia, "../.devprivops/tests/a/out.dpia.yml")

    # ASVS

    dt_2_1 = dt1.clone(categories=[
        "dpia:confidential",
        "dpia:personal",
        "dpia:authenticated only",
        "dpia:sensitive"
    ])

    ds_2_1 = ds1.clone(
        environment=["browser"]
    )

    dfd2 = dfd.clone(data_types=[dt_2_1, dt2, dt3], data_stores=[ds_2_1])
    dpia2 = dpia.clone()

    to_yaml(dfd2, "../.devprivops/tests/asvs_browser/out.dfd.yml")
    to_yaml(dpia2, "../.devprivops/tests/asvs_browser/out.dpia.yml")

    # DPIA con

    pd_3_1 = pd1.clone(
        destinataries=[":dont exist", ":external system", *pd1.destinataries],
        retention_period="wrong",
        # bla=1,
    )
    dt_3_1 = DataType(
        "new data type",
        [],
        "1m",
        ["dpia:personal"]
    )

    proc_3_1 = Process(
        "new proc",
        [],
        [],
        [],
        [],
        ["new purpose 1"],
        [],
        [],
    )

    processing_3_1 = processing1.clone(risks=["new risk"])

    dfd3 = dfd.clone(
        data_types=[dt_3_1, *dfd.data_types],
        processes=[proc_3_1, *dfd.processes],
    )
    dpia3 = dpia.clone(
        personal_data=[pd_3_1, pd2, pd3],
        personal_data_processing_whitelist=["new purpose 2"],
        personal_data_processing_that_requires_DPIA=["new purpose 3"],
        processings=[processing_3_1],
        last_update="wrong date",
    )

    to_yaml(dfd3, "../.devprivops/tests/dpia_con/out.dfd.yml")
    to_yaml(dpia3, "../.devprivops/tests/dpia_con/out.dpia.yml")

    data_4_1 = DataStored(
        ":type 3",
        "eternal",
        ":C new",
        ":R new",
        ":U new",
        ":D new",
    )

    ds_4_1 = DataStore(
        "new db",
        [data_4_1],
        ["Portugal"],
        [],
        [],
        []
    )

    ee_4_1 = ee1.clone(
        id_="new ent",
        consumes=[":type 1"],
        produces=[":type 2"],
    )

    df_4_1 = df1.clone(
        from_=":message",
        to=":no exists",
        data=[":type 4", *df1.data],
        periodicity="invalid",
    )

    dt_4_1 = dt1.clone(
        aggregates=[":type 5", *dt1.aggregates],
        validity="invalid",
    )

    processing_4_1 = processing1.clone(
        legitimate_interest=[":human", ":User", ":new one"],
        vital_interest=[":human", ":User", ":another one"],
    )

    dfd4 = dfd.clone(
        external_entities=[ee_4_1, *dfd.external_entities],
        data_stores=[ds_4_1, *dfd.data_stores],
        data_flows=[df_4_1, *dfd.data_flows],
        data_types=[dt_4_1, *dfd.data_types],
    )
    dpia4 = dpia.clone(
        processings=[processing_4_1],
    )

    to_yaml(dfd4, "../.devprivops/tests/gdpr_con/out.dfd.yml")
    to_yaml(dpia4, "../.devprivops/tests/gdpr_con/out.dpia.yml")
