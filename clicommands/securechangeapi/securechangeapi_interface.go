package securechangeapi

import (
	"uqtu/mediator/console"
	"uqtu/mediator/scworkflow"
)

type SecurchangeAPIManager interface {
	GetSecurechangeWorkflowTriggers() (*scworkflow.WorkflowTriggers, error)
	CreateSecurechangeWorkflowTriggers(wf_triggers *scworkflow.WorkflowTriggers) error
	DeleteSecurechangeWorkflowTriggers(wf_trigger_id int) error
	GetSecurechangeWorkflowTriggerByID(id int) (*scworkflow.WorkflowTrigger, error)
	GetSecurechangeWorkflows(get_steps bool) (*scworkflow.Workflows, error)
}

//local implementation of the interface

type scManager struct {
	SC_username string
	SC_pwd      string
	SC_host     string
}

func (mgr *scManager) GetSecurechangeWorkflowTriggers() (*scworkflow.WorkflowTriggers, error) {
	if err := mgr.getSecurechangeCredentials(); err != nil {
		return nil, err
	}
	return scworkflow.GetSecurechangeWorkflowTriggers(mgr.SC_username, mgr.SC_pwd, mgr.SC_host)
}

func (mgr *scManager) CreateSecurechangeWorkflowTriggers(wf_triggers *scworkflow.WorkflowTriggers) error {
	if err := mgr.getSecurechangeCredentials(); err != nil {
		return err
	}
	return scworkflow.CreateSecurechangeWorkflowTriggers(wf_triggers, mgr.SC_username, mgr.SC_pwd, mgr.SC_host)
}

func (mgr *scManager) DeleteSecurechangeWorkflowTriggers(id int) error {
	if err := mgr.getSecurechangeCredentials(); err != nil {
		return err
	}
	return scworkflow.DeleteSecurechangeWorkflowTriggers(id, mgr.SC_username, mgr.SC_pwd, mgr.SC_host)
}

func (mgr *scManager) GetSecurechangeWorkflowTriggerByID(id int) (*scworkflow.WorkflowTrigger, error) {
	if err := mgr.getSecurechangeCredentials(); err != nil {
		return nil, err
	}
	return scworkflow.GetSecurechangeWorkflowTriggerByID(id, mgr.SC_username, mgr.SC_pwd, mgr.SC_host)
}

func (mgr *scManager) GetSecurechangeWorkflows(get_steps bool) (*scworkflow.Workflows, error) {
	if err := mgr.getSecurechangeCredentials(); err != nil {
		return nil, err
	}
	return scworkflow.GetSecurechangeWorkflows(mgr.SC_username, mgr.SC_pwd, mgr.SC_host, get_steps)
}

func (mgr *scManager) getSecurechangeCredentials() error {
	var err error
	// ask for Securechange host if not provided via dedicated flag
	for {
		if mgr.SC_host != "" {
			break
		}
		if mgr.SC_host, err = console.GetText("SecureChange Host"); err != nil {
			return err
		}
	}
	// ask for user name if not provided via dedicated flag
	for {
		if mgr.SC_username != "" {
			break
		}
		if mgr.SC_username, err = console.GetText("SecureChange Username"); err != nil {
			return err
		}
	}

	// ask for password if not provided via dedicated flag
	for {
		if mgr.SC_pwd != "" {
			break
		}
		if mgr.SC_pwd, err = console.GetPassword("SecureChange Password"); err != nil {
			return err
		}
	}

	return nil
}
