import json
from typing import List

import yaml

from dfd import DataType, ExternalEntity, Process, DataStore, DataStored, DataFlow, DFD, DataFormat, Data, \
    ProcessedCategorization
from dpia import DPO, PersonalDatum, DPIA, Risk, Processing, Purpose, SupervisoryAuthorityVeredict


def to_yaml(data, fname):
    data = json.dumps(data, default=lambda x: {
        k.replace("_", " ").rstrip(): x.__dict__[k]
        for k in x.__dict__
    })

    with open(fname, 'w') as f:
        f.write(yaml.dump(json.loads(data)))


def to_data_list(dt_names: List[str]) -> List[Data]:
    return [Data(e, DataFormat.default(), ProcessedCategorization.default()) for e in dt_names]

# yaml-language-server: $schema=../../schemas/dfd-schema.json
if __name__ == '__main__':
    # Data(e, DataFormat("psdc:plain", "plain text", "business")

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
        to_data_list([":message"]),
        to_data_list([":message"]),
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
        to_data_list([":message"]),
        to_data_list([":message"]),
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
        to_data_list([":message"]),
        to_data_list([":message"]),
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
        DataFormat.default(),
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
        to_data_list([":message"]),
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
        to_data_list([":message"]),
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
        1,
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
        scores_users=False,
        automated_decisions=False,
        legal_impact_for_the_user=False,
        systematic_monitoring=False,
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
            # [],
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

    to_yaml(dfd, "../examples/global/tests/a/out.dfd.yml")
    to_yaml(dpia, "../examples/global/tests/a/out.dpia.yml")

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

    to_yaml(dfd2, "../examples/global/tests/asvs_browser/out.dfd.yml")
    to_yaml(dpia2, "../examples/global/tests/asvs_browser/out.dpia.yml")

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

    to_yaml(dfd3, "../examples/global/tests/dpia_con/out.dfd.yml")
    to_yaml(dpia3, "../examples/global/tests/dpia_con/out.dpia.yml")

    data_4_1 = DataStored(
        ":type 3",
        "eternal",
        ":C new",
        ":R new",
        ":U new",
        ":D new",
        DataFormat.default()
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
        consumes=to_data_list([":type 1"]),
        produces=to_data_list([":type 2"]),
    )

    df_4_1 = df1.clone(
        from_=":message",
        to=":no exists",
        data=[Data(":type 4", DataFormat.default(), ProcessedCategorization.default()), *df1.data],
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

    to_yaml(dfd4, "../examples/global/tests/gdpr_con/out.dfd.yml")
    to_yaml(dpia4, "../examples/global/tests/gdpr_con/out.dpia.yml")

    risk_5_1 = risk1.clone(
        impact=10,
        likelyhood=10,
    )

    purpose_5_1 = purpose1.clone(
        adequate=False,
        relevant=False,
        limited=False,
    )

    processing_5_1 = processing1.clone(
        requires_new_technologies=True,
        scores_users=True,
        automated_decisions=True,
        legal_impact_for_the_user=True,
        systematic_monitoring=True,
        large_scale_processing=True,
        supervisory_authority_veredict=SupervisoryAuthorityVeredict([], True),
        lawful=False,
        fair=False,
        transparent=False,
        specific=False,
        explicit=False,
        legitimate=False,
        purposes=[purpose_5_1]
    )

    proc_5_1 = proc1.clone(
        purposes=["dpia:new purpose", *proc1.purposes]
    )

    dfd5 = dfd.clone(
        processes=[proc_5_1, *dfd.processes[1:]]
    )
    dpia5 = dpia.clone(
        processings=[processing_5_1],
        personal_data_processing_whitelist=[],
        # personal_data_processing_that_requires_DPIA=[":message routing", *dpia.personal_data_processing_that_requires_DPIA]
        risks=[risk_5_1],
    )

    to_yaml(dfd5, "../examples/global/tests/dpia_pol/out.dfd.yml")
    to_yaml(dpia5, "../examples/global/tests/dpia_pol/out.dpia.yml")

    ee_6_1 = ee1.clone(
        location=["Nowhere"],
    )

    proc_6_1 = proc1.clone(
        location=["Farlands"],
        create=[":message"],
    )

    df_6_1 = df1.clone(
        data=[Data(":new data", DataFormat.default(), ProcessedCategorization.default()), *df1.data]
    )

    data_6_1 = data1.clone(
        create=":message",
    )

    ds_6_1 = ds1.clone(
        location=["Don't exist"],
        data_stored=[data_6_1],
    )

    dfd6 = dfd.clone(
        external_entities=[ee_6_1, *dfd.external_entities[1:]],
        data_stores=[ds_6_1, *dfd.data_stores[1:]],
        processes=[proc_6_1, *dfd.processes[1:]],
        data_flows=[df_6_1, *dfd.data_flows[1:]]
    )
    dpia6 = dpia.clone()

    to_yaml(dfd6, "../examples/global/tests/dfd_pol/out.dfd.yml")
    to_yaml(dpia6, "../examples/global/tests/dfd_pol/out.dpia.yml")

    dt_p_1 = DataType(
        "message",
        [],
        "eternal",
        ["dpia:confidential", "dpia:personal"],
    )
    dt_p_2 = DataType(
        "AccountId",
        [],
        "eternal",
        ["dpia:personal"],
    )
    dt_p_3 = DataType(
        "Account",
        [],
        "eternal",
        ["dpia:confidential", "dpia:personal"],
    )

    ee_p_1 = ExternalEntity(
        "dpia:User",
        [Data(":message", DataFormat.default(), ProcessedCategorization.default())],
        [Data(":message", DataFormat.default(), ProcessedCategorization.default())],
        ["Portugal"],
        [],
        ["dpia:human"],
        ">16",
        False,
        [],
        [],
    )
    ee_p_2 = ExternalEntity(
        "dpia:AccountSystem",
        [Data(":message", DataFormat.default(), ProcessedCategorization.default())],
        [Data(":message", DataFormat.default(), ProcessedCategorization.default())],
        ["Portugal"],
        [],
        ["dpia:human"],
        ">16",
        False,
        [],
        [],
    )

    proc_p_1 = Process(
        "send message",
        [Data(":message", DataFormat.default(), ProcessedCategorization.default())],
        [Data(":message", DataFormat.default(), ProcessedCategorization.default())],
        ["Portugal"],
        [],
        ["dpia:message routing"],
        [],
        [],
    )

    ds_p_1 = DataStore(
        "message db",
        [
            DataStored(
                ":message",
                "eternal",
                ":C store message",
                ":R store message",
                ":U store message",
                ":D store message",
                DataFormat.default(),
            )
        ],
        ["Portugal"],
        [],
        [],
        [],
    )

    store_message_df_p = [
        DataFlow(
            f"{m} store message",
            ":send message",
            ":message db",
            [Data(":message", DataFormat.default(), ProcessedCategorization.default())],
            "signal",
            "1m",
            1,
            [],
            [],
        )
        for m in ('C', 'R', 'U', 'D')
    ]
    message_df_p = [
        DataFlow(
            f"{m} message",
            "dpia:User",
            ":send message",
            [Data(":message", DataFormat.default(), ProcessedCategorization.default())],
            "signal",
            "1m",
            1,
            [],
            [],
        )
        for m in ('C', 'R', 'U', 'D')
    ]
    deliver_df_p = DataFlow(
            f"Deliver message",
            ":send message",
            "dpia:User",
            [Data(":message", DataFormat.default(), ProcessedCategorization.default())],
            "signal",
            "1m",
            1,
            [],
            [],
        )

    dpo_p = [
        DPO(":The", "the@email.com"),
        DPO(":Man", "man@email.com"),
    ]

    pd_p_1 = PersonalDatum(
        "dfd:message",
        ":personal",
        [],
        [],
        [":User"],
        "2d",
        False,
        [":message routing"],
        [],
    )
    pd_p_2 = PersonalDatum(
        "dfd:AccountId",
        ":personal",
        [],
        [],
        [":User"],
        "2d",
        False,
        [],
        [],
    )
    pd_p_3 = PersonalDatum(
        "dfd:Account",
        ":personal",
        [],
        [],
        [":User"],
        "2d",
        False,
        [],
        [],
    )

    r_p_1 = Risk(
        "Risk 1",
        0,
        1,
        []
    )

    processing_p_1 = Processing(
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
        scores_users=False,
        automated_decisions=False,
        legal_impact_for_the_user=False,
        systematic_monitoring=False,
        large_scale_processing=False,
        lawful=True,
        fair=True,
        transparent=True,
        specific=True,
        explicit=True,
        legitimate=True,
        purposes=[purpose1],
        risks=[f":{r_p_1.id_}"],
        supervisory_authority_veredict=SupervisoryAuthorityVeredict(
            [":Supervisor"],
            True
        )
    )

    dfd_perfect = DFD(
        [dt_p_1, dt_p_2, dt_p_3],
        [ee_p_1, ee_p_2],
        [proc_p_1],
        [ds_p_1],
        [*store_message_df_p, *message_df_p, deliver_df_p]
    )
    dpia_perfect = DPIA(
        ":last_update",
        [":Someone", ":Else"],
        dpo_p,
        [pd_p_1, pd_p_2, pd_p_3],
        [r_p_1],
        [":message routing"],
        [],
        [processing_p_1]
    )

    to_yaml(dfd_perfect, "../examples/global/tests/perfect_scenario/perfect_scenario.dfd.yml")
    to_yaml(dpia_perfect, "../examples/global/tests/perfect_scenario/perfect_scenario.dpia.yml")
