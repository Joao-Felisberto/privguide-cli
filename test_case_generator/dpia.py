from test_case_generator.util import Cloneable


class Risk(Cloneable):
    def __init__(
            self,
            id_: str,
            impact: int,
            likelyhood: int,
            accepted_mitigations: [str],
    ):
        self.id_: str = id_
        self.impact: int = impact
        self.likelyhood: int = likelyhood
        self.accepted_mitigations: [str] = accepted_mitigations


class SupervisoryAuthorityVeredict(Cloneable):
    def __init__(
            self,
            contact: [str],
            allowed: bool,
    ):
        self.contact: [str] = contact
        self.allowed: bool = allowed


class Purpose(Cloneable):
    def __init__(
            self,
            id_: str,
            adequate: bool,
            relevant: bool,
            limited: bool,
    ):
        self.id_: str = id_
        self.adequate: bool = adequate
        self.relevant: bool = relevant
        self.limited: bool = limited


class Processing(Cloneable):
    def __init__(
            self,
            id_: str,
            automated_decisions: bool,
            explicit: bool,
            fair: bool,
            is_official_authority: bool,
            large_scale_processing: bool,
            lawful: bool,
            legal_impact_for_the_user: bool,
            legally_mandated: bool,
            legitimate: bool,
            legitimate_interest: [str],
            professional_secrecy: bool,
            public_interest: bool,
            purposes: [Purpose],
            required_for_contract: [str],
            requires_new_technologies: bool,
            risk_to_rights_and_freedoms_of: [str],
            risks: [str],
            scores_users: bool,
            specific: bool,
            supervisory_authority_veredict: SupervisoryAuthorityVeredict,
            systematic_monitoring: bool,
            transparent: bool,
            vital_interest: [str],
    ):
        self.id_: str = id_
        self.automated_decisions: bool = automated_decisions
        self.explicit: bool = explicit
        self.fair: bool = fair
        self.is_official_authority: bool = is_official_authority
        self.large_scale_processing: bool = large_scale_processing
        self.lawful: bool = lawful
        self.legal_impact_for_the_user: bool = legal_impact_for_the_user
        self.legally_mandated: bool = legally_mandated
        self.legitimate: bool = legitimate
        self.legitimate_interest: [str] = legitimate_interest
        self.professional_secrecy: bool = professional_secrecy
        self.public_interest: bool = public_interest
        self.purposes: [Purpose] = purposes
        self.required_for_contract: [str] = required_for_contract
        self.requires_new_technologies: bool = requires_new_technologies
        self.risk_to_rights_and_freedoms_of: [str] = risk_to_rights_and_freedoms_of
        self.risks: [str] = risks
        self.scores_users: bool = scores_users
        self.specific: bool = specific
        self.supervisory_authority_veredict: SupervisoryAuthorityVeredict = supervisory_authority_veredict
        self.systematic_monitoring: bool = systematic_monitoring
        self.transparent: bool = transparent
        self.vital_interest: [str] = vital_interest


class PersonalDatum(Cloneable):
    def __init__(
            self,
            id_: str,
            kind: str,
            required_by_law: [str],
            necessary_to_enter_contract: [str],
            destinataries: [str],
            retention_period: str,
            abides_by_code_of_conduct: bool,
            purposes: [str],
            transfers_to_third_parties: [str],
    ):
        self.id_: str = id_
        self.kind: str = kind
        self.required_by_law: [str] = required_by_law
        self.necessary_to_enter_contract: [str] = necessary_to_enter_contract
        self.destinataries: [str] = destinataries
        self.retention_period: str = retention_period
        self.abides_by_code_of_conduct: bool = abides_by_code_of_conduct
        self.purposes: [str] = purposes
        self.transfers_to_third_parties: [str] = transfers_to_third_parties


class DPO(Cloneable):
    def __init__(
            self,
            name: str,
            contact: str,
    ):
        self.name: str = name
        self.contact: str = contact


class DPIA(Cloneable):
    def __init__(
            self,
            last_update: str,
            responsible: [str],
            DPO_: [DPO],
            personal_data: [PersonalDatum],
            risks: [Risk],
            personal_data_processing_whitelist: [str],
            personal_data_processing_that_requires_DPIA: [str],
            processings: [Processing],
    ):
        self.last_update: str = last_update
        self.responsible: [str] = responsible
        self.DPO_: [DPO] = DPO_
        self.personal_data: [PersonalDatum] = personal_data
        self.risks: [Risk] = risks
        self.personal_data_processing_whitelist: [str] = personal_data_processing_whitelist
        self.personal_data_processing_that_requires_DPIA: [str] = personal_data_processing_that_requires_DPIA
        self.processings: [Processing] = processings
